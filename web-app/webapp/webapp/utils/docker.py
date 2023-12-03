from dataclasses import dataclass
from pathlib import Path

import docker.errors
import tomlkit
from docker.models.containers import Container


@dataclass
class ServiceConfig:
    name: str
    image: str
    db_host: str
    db_port: int
    db_name: str
    db_user: str
    db_password: str

    def to_dict(self):
        return {
            "name": self.name,
            "image": self.image,
            "db_host": self.db_host,
            "db_port": self.db_port,
            "db_name": self.db_name,
            "db_user": self.db_user,
            "db_password": self.db_password,
        }


@dataclass
class LoadBalancerConfig:
    image: str
    strategy: str

    def to_dict(self):
        return {
            "image": self.image,
            "strategy": self.strategy,
        }


@dataclass
class AutoScalerConfig:
    image: str
    strategy: str
    cpu_threshold: float
    mem_threshold: float
    network: str
    service_config: ServiceConfig
    lb_config: LoadBalancerConfig

    def to_dict(self):
        return {
            "strategy": self.strategy,
            "cpu_threshold": self.cpu_threshold,
            "mem_threshold": self.mem_threshold,
            "network": self.network,
            "service": self.service_config.to_dict(),
            "load_balancer": self.lb_config.to_dict(),
        }


@dataclass
class MySQLConfig:
    root_password: str
    db_name: str
    db_user: str
    db_password: str
    sql_file: Path


def create_network(name: str):
    docker_cli = docker.from_env()
    docker_cli.networks.create(name)


def create_auto_scaler(name: str, conf_file_path: Path, config: AutoScalerConfig):
    with conf_file_path.open("w") as f:
        tomlkit.dump(config.to_dict(), f)

    client = docker.from_env()

    client.images.pull(config.image)

    container: Container = client.containers.create(
        config.image,
        detach=True,
        hostname=name,
        init=True,
        name=name,
        network=config.network,
        restart_policy={"Name": "always"},
        volumes=[f"{conf_file_path}:/config.toml:ro", "/var/run/docker.sock:/var/run/docker.sock"],
    )
    container.start()


def create_mysql_container(image: str, name: str, network: str, config: MySQLConfig) -> str:
    client = docker.from_env()

    vol_name = f"data-{name}"
    try:
        client.volumes.create(vol_name)
    except docker.errors.APIError:
        ...

    client.images.pull(image)
    try:
        container = client.containers.create(
            image,
            detach=True,
            environment={
                "MYSQL_ROOT_PASSWORD": config.root_password,
                "MYSQL_DATABASE": config.db_name,
                "MYSQL_USER": config.db_user,
                "MYSQL_PASSWORD": config.db_password,
            },
            hostname=name,
            name=name,
            network=network,
            restart_policy={"Name": "always"},
            volumes=[f"{config.sql_file}:/docker-entrypoint-initdb.d", f"{vol_name}:/var/lib/mysql"],
        )
    except docker.errors.APIError:
        container: docker.models.containers.Container = client.containers.get(name)

    container.start()

    return container.attrs["NetworkSettings"]["Networks"][network]["IPAddress"]
