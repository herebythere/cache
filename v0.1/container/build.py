import argparse
import json
import os
import subprocess

from string import Template

# constants
trailing_slash = "/"

store_data_dir = "data/"

container_compose_template_filepath = "podman-compose.yml.template"
container_compose_filepath = "podman-compose.yml"
conf_template_filepath = "redis.conf.template"
conf_filepath = "redis.conf"
container_template_filepath = "redis.podmanfile.template"
container_filepath = "redis.podmanfile"

# defaults

default_config_filepath = "config/cache.json"
default_dest_filepath = "cache/"
default_template_filepath = "templates/"

# parser

parser = argparse.ArgumentParser()
parser.add_argument(
    "--dest", help="provide a preferred desitnation for build results")
parser.add_argument("--templates", help="preferred template directory")
parser.add_argument(
    "--config", help="override everything with a json config file")


def apply_defaults_to_args(args):
    if args["config"] == None:
        args["config"] = default_config_filepath
    if args["templates"] == None:
        args["templates"] = default_template_filepath
    if args["dest"] == None:
        args["dest"] = default_dest_filepath

    return args


def get_filepaths(args):
    filepaths = {
        "container_compose_template": args["templates"] + container_compose_template_filepath,
        "container_compose": args["dest"] + container_compose_filepath,
        "conf_template": args["templates"] + conf_template_filepath,
        "conf": args["dest"] + conf_filepath,
        "container_template": args["templates"] + container_template_filepath,
        "container": args["dest"] + container_filepath,
    }

    return filepaths


def get_config(source):
    config_file = open(source, "r")
    config = json.load(config_file)
    config_file.close()

    return config


def create_required_directories(args):
    data_dir = args["dest"] + store_data_dir
    if not os.path.exists(data_dir):
        os.makedirs(data_dir)


def create_template(source, target, keywords):
    source_file = open(source, "r")
    source_file_template = Template(source_file.read())
    source_file.close()

    updated_source_file_template = source_file_template.substitute(**keywords)

    target_file = open(target, "w+")
    target_file.write(updated_source_file_template)
    target_file.close()


def create_required_templates(args, filepaths, config):
    args_and_config_map = {}
    args_and_config_map.update(args)
    args_and_config_map.update(config)

    create_template(filepaths["container_compose_template"],
                    filepaths["container_compose"],
                    args_and_config_map)

    create_template(filepaths["conf_template"],
                    filepaths["conf"],
                    args_and_config_map)

    create_template(filepaths["container_template"],
                    filepaths["container"],
                    args_and_config_map)


def compose_cache_with_podman(filepaths):
    subprocess.run(["podman-compose", "--file",
                   filepaths['container_compose'], "build"])


def build_cache_with_podman(args, filepaths, config):
    create_required_directories(args)
    create_required_templates(args, filepaths, config)
    compose_cache_with_podman(filepaths)


if __name__ == "__main__":
    args = vars(parser.parse_args())
    args = apply_defaults_to_args(args)
    filepaths = get_filepaths(args)
    config = get_config(args["config"])

    build_cache_with_podman(args, filepaths, config)
