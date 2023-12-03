from django.urls import path

from webapp.deployment.views import deployment_detail_view
from webapp.deployment.views import deployment_list_view
from webapp.deployment.views import deployment_update_view

app_name = "deployment"
urlpatterns = [
    path("", view=deployment_list_view, name="list"),
    path("<int:pk>/", view=deployment_detail_view, name="detail"),
    path("new/", view=deployment_update_view, name="new"),
    path("update/<int:pk>/", view=deployment_update_view, name="update"),
]
