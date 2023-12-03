from django import forms

from webapp.deployment.models import Deployment


class DeploymentForm(forms.ModelForm):
    def __init__(self, *args, **kwargs):
        self.user = kwargs.pop("user")
        super().__init__(*args, **kwargs)

    def save(self, commit: bool = True):
        obj = super().save(commit=False)
        obj.user = self.user
        if commit:
            obj.save()
        return obj

    class Meta:
        model = Deployment
        fields = ("service", "lb_strategy", "scaler_strategy", "scaler_cpu_threshold", "scaler_mem_threshold")
