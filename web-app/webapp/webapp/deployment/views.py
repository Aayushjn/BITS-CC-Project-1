import hashlib

from django.conf import settings
from django.contrib.auth.mixins import LoginRequiredMixin
from django.contrib.messages.views import SuccessMessageMixin
from django.urls import reverse
from django.utils.translation import gettext_lazy as _
from django.views.generic import DetailView
from django.views.generic import ListView
from django.views.generic import UpdateView

import webapp.utils.docker as docker_utils
from webapp.deployment.forms import DeploymentForm
from webapp.deployment.models import Deployment


class DeploymentDetailView(LoginRequiredMixin, DetailView):
    model = Deployment
    slug_field = "id"
    slug_url_kwarg = "id"


deployment_detail_view = DeploymentDetailView.as_view()


class DeploymentListView(LoginRequiredMixin, ListView):
    model = Deployment

    def get_queryset(self):
        return super().get_queryset().filter(user=self.request.user)


deployment_list_view = DeploymentListView.as_view()


class DeploymentUpdateView(LoginRequiredMixin, SuccessMessageMixin, UpdateView):
    model = Deployment
    success_message = _("Information successfully updated")
    form_class = DeploymentForm

    def get_object(self, queryset=None):
        try:
            return super().get_object(queryset)
        except AttributeError:
            return None

    def get_success_url(self):
        return reverse("deployment:list")

    def get_form_kwargs(self):
        kwargs = super().get_form_kwargs()
        kwargs["user"] = self.request.user
        return kwargs

    def form_valid(self, form):
        resp = super().form_valid(form)

        user_hash = hashlib.sha1(self.request.user.name.encode("utf-8")).hexdigest()
        docker_utils.create_network(name=f"net-{user_hash}")

        mysql_config = docker_utils.MySQLConfig(
            root_password=settings.MYSQL_ROOT_PASSWORD,
            db_name=settings.MYSQL_DATABASE,
            db_user=settings.MYSQL_USER,
            db_password=settings.MYSQL_PASSWORD,
            sql_file=settings.SQL_FILE,
        )

        ip_addr = docker_utils.create_mysql_container(
            image=settings.MYSQL_IMAGE,
            name=f"mysql-{user_hash}",
            network=f"net-{user_hash}",
            config=mysql_config,
        )

        config = docker_utils.AutoScalerConfig(
            image=f"{settings.DOCKER_REPO}/{settings.SCALER_IMAGE}",
            strategy=form.cleaned_data["scaler_strategy"],
            cpu_threshold=form.cleaned_data["scaler_cpu_threshold"],
            mem_threshold=form.cleaned_data["scaler_mem_threshold"],
            network=f"net-{user_hash}",
            service_config=docker_utils.ServiceConfig(
                name=form.cleaned_data["service"],
                image=f"{settings.DOCKER_REPO}/{settings.SERVICE_PREFIX}{form.cleaned_data['service']}",
                db_host=ip_addr,
                db_port=3306,
                db_name=settings.MYSQL_DATABASE,
                db_user=settings.MYSQL_USER,
                db_password=settings.MYSQL_PASSWORD,
            ),
            lb_config=docker_utils.LoadBalancerConfig(
                image=f"{settings.DOCKER_REPO}/{settings.LB_IMAGE}",
                strategy=form.cleaned_data["lb_strategy"],
            ),
        )

        docker_utils.create_auto_scaler(
            name=f"scaler-{user_hash}-{form.cleaned_data['service']}",
            conf_file_path=settings.SERVICE_CONF_DIR.joinpath(f"scaler-{user_hash}-{form.cleaned_data['service']}.toml"),
            config=config,
        )

        return resp


deployment_update_view = DeploymentUpdateView.as_view()
