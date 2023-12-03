from django.apps import AppConfig
from django.utils.translation import gettext_lazy as _


class DeploymentConfig(AppConfig):
    name = "webapp.deployment"
    verbose_name = _("Deployments")
    default_auto_field = "django.db.models.BigAutoField"
