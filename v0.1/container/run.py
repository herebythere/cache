import argparse
import subprocess

default_config_filepath = "./dest/podman-compose.yml"

parser = argparse.ArgumentParser()
parser.add_argument(
    "--file", help="provide an alternative podman-compose file")


def apply_defaults_to_args(args):
    if args["file"] == None:
        args["file"] = default_config_filepath

    return args


def run_cache_with_podman(args):
    subprocess.run(["podman-compose", "--file",
                   args["file"], "up", "--detach"])


if __name__ == "__main__":
    args = vars(parser.parse_args())
    args = apply_defaults_to_args(args)

    run_cache_with_podman(args)
