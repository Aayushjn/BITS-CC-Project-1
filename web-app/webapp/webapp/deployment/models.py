from django.contrib.auth import get_user_model
from django.db import models
from django.utils.translation import gettext_lazy as _


class LoadBalancerStrategy(models.TextChoices):
    ROUND_ROBIN = "round_robin", _("Round Robin")
    LEAST_CONNECTIONS = "least_conns", _("Least Connections")
    POWER_OF_2 = "power_of_2", _("Power of 2")
    RANDOM = "random", _("Random")
    IP_HASH = "ip_hash", _("IP Hash")
    WEIGHTED_ROUND_ROBIN = "weighted_round_robin", _("Weighted Round Robin")


class AutoScalerStrategy(models.TextChoices):
    THRESHOLD = "threshold", _("Threshold-based")
    TIME_SERIES = "time_series", _("Time Series-based")


class ServiceType(models.TextChoices):
    BOOKING = "booking", _("Booking")
    HOTEL = "hotel", _("Hotel")
    PACKAGES = "packages", _("Packages")
    TICKET = "ticket", _("Ticket")
    USER = "user", _("User")


class Deployment(models.Model):
    user = models.ForeignKey(get_user_model(), on_delete=models.CASCADE)
    service = models.CharField(max_length=100, choices=ServiceType.choices)
    lb_strategy = models.CharField(max_length=100, choices=LoadBalancerStrategy.choices)
    scaler_strategy = models.CharField(max_length=100, choices=AutoScalerStrategy.choices)
    scaler_cpu_threshold = models.DecimalField(max_digits=5, decimal_places=2)
    scaler_mem_threshold = models.DecimalField(max_digits=5, decimal_places=2)

    def __str__(self):
        return self.service

    def __repr__(self):
        return (
            f"Deployment({self.service=}, {self.lb_strategy=}, {self.scaler_strategy=}, {self.scaler_cpu_threshold=}, "
            f"{self.scaler_mem_threshold=})"
        )

    def get_absolute_url(self):
        from django.urls import reverse

        return reverse("deployment:detail", kwargs={"pk": self.pk})
