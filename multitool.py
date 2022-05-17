#!/usr/bin/env python

import sys
import requests
import json
import click
from loguru import logger


class MatrixPostRequestFailed(Exception):
    pass


class MatrixAPI:
    def __init__(self, url: str, **kwargs):
        api_prefix = kwargs.get("api_prefix", "/_matrix/client/r0")
        self.base = url.strip().rstrip("/") + api_prefix

    def post_json(self, endpoint: str, data: dict):
        endpoint = endpoint.lstrip("/")
        url = self.base + "/" + endpoint

        logger.debug(f"json post to {url}")
        r = requests.post(url, json=data)
        if r.status_code != 200:
            logger.error("response was: " + r.content.decode("utf8"))
            raise MatrixPostRequestFailed(f"status code {r.status_code}")

        return r.json()
        
    def register(self, username: str, password: str):
        data = self.post_json(endpoint="register", data={
            "username": username,
            "password": password,
            "auth": {
                "type": "m.login.dummy",
            },
        })

        logger.info(f"registered user {username} successfully")
        print("response: " + json.dumps(data, indent=2))

    def login(self, username: str, password: str):
        data = self.post_json(endpoint="login", data={
            "identifier": {
                "type": "m.id.user",
                "user": username,
                },
            "password": password,
            "type": "m.login.password",
        })

        logger.info("login successful")
        print("response: " + json.dumps(data, indent=2))


@click.group()
@click.option('--debug/--no-debug', default=False)
def cli(debug):
    logger.remove()
    if debug:
        logger.add(sys.stderr, level="DEBUG")
    else:
        logger.add(sys.stderr, level="INFO")


@cli.command()
@click.option("--url", required=True)
@click.argument("username")
@click.argument("password")
def register(url, username, password):
    m = MatrixAPI(url)
    m.register(username, password)


@cli.command()
@click.option("--url", required=True)
@click.argument("username")
@click.argument("password")
def login(url, username, password):
    m = MatrixAPI(url)
    m.login(username, password)


if __name__ == "__main__":
    cli()
