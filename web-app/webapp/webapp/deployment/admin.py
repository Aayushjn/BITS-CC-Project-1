from django.contrib import admin

from webapp.deployment.models import Deployment


@admin.register(Deployment)
class UserAdmin(admin.ModelAdmin):
    pass
