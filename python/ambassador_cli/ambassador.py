# Copyright 2018-2020 Datawire. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

########
# This is the ambassador CLI. Despite the impression given by its name, it is actually
# primarily a debugging tool at this point: the most useful thing to do with it is to
# run "ambassador dump --watt path-to-watt-snapshot-file" and have it spit out the IR,
# etc.
########

from typing import ClassVar, Optional, Set, TYPE_CHECKING
from typing import cast as typecast

import sys

import cProfile
import json
import logging
import os
import orjson
import pstats
import signal
import traceback

import click

from ambassador import Scout, Config, IR, Diagnostics, Version
from ambassador.fetch import ResourceFetcher
from ambassador.envoy import EnvoyConfig, V3Config

from ambassador.utils import (
    RichStatus,
    SecretHandler,
    SecretInfo,
    NullSecretHandler,
    Timer,
    parse_json,
    dump_json,
)

if TYPE_CHECKING:
    from ambassador.ir import IRResource  # pragma: no cover

__version__ = Version

logging.basicConfig(
    level=logging.INFO,
    format="%%(asctime)s ambassador-cli %s %%(levelname)s: %%(message)s" % __version__,
    datefmt="%Y-%m-%d %H:%M:%S",
)

logger = logging.getLogger("ambassador")


def handle_exception(what, e, **kwargs):
    tb = "\n".join(traceback.format_exception(*sys.exc_info()))

    scout = Scout()
    result = scout.report(action=what, mode="cli", exception=str(e), traceback=tb, **kwargs)

    logger.debug("Scout %s, result: %s" % ("enabled" if scout._scout else "disabled", result))

    logger.error("%s: %s\n%s" % (what, e, tb))

    show_notices(result)


def show_notices(result: dict, printer=logger.log):
    notices = result.get("notices", [])

    for notice in notices:
        lvl = logging.getLevelName(notice.get("level", "ERROR"))

        printer(lvl, notice.get("message", "?????"))


def stdout_printer(lvl, msg):
    print("%s: %s" % (logging.getLevelName(lvl), msg))


def version():
    """
    Show Ambassador's version
    """

    print("Ambassador %s" % __version__)

    scout = Scout()

    print("Ambassador Scout version %s" % scout.version)
    print("Ambassador Scout semver  %s" % scout.get_semver(scout.version))

    result = scout.report(action="version", mode="cli")
    show_notices(result, printer=stdout_printer)


def showid():
    """
    Show Ambassador's installation ID
    """

    scout = Scout()

    print("Ambassador Scout installation ID %s" % scout.install_id)

    result = scout.report(action="showid", mode="cli")
    show_notices(result, printer=stdout_printer)


def file_checker(path: str) -> bool:
    logger.debug("CLI file checker: pretending %s exists" % path)
    return True


class CLISecretHandler(SecretHandler):
    # HOOK: if you're using dump and you need it to pretend that certain missing secrets
    # are present, add them to LoadableSecrets. At Some Point(tm) there will be a switch
    # to add these from the command line, but Flynn didn't actually need that for the
    # debugging he was doing...

    LoadableSecrets: ClassVar[Set[str]] = set(
        # "ssl-certificate.mynamespace"
    )

    def load_secret(
        self, resource: "IRResource", secret_name: str, namespace: str
    ) -> Optional[SecretInfo]:
        # Only allow a secret to be _loaded_ if it's marked Loadable.

        key = f"{secret_name}.{namespace}"

        if key in CLISecretHandler.LoadableSecrets:
            self.logger.info(f"CLISecretHandler: loading {key}")
            return SecretInfo(
                secret_name,
                namespace,
                "mocked-loadable-secret",
                "-mocked-cert-",
                "-mocked-key-",
                decode_b64=False,
            )

        self.logger.debug(f"CLISecretHandler: cannot load {key}")
        return None


@click.command()
@click.argument("config_dir_path", type=click.Path())
@click.option("--secret-dir-path", type=click.Path(), help="Directory into which to save secrets")
@click.option("--watt", is_flag=True, help="If set, input must be a WATT snapshot")
@click.option("--debug", is_flag=True, help="If set, generate debugging output")
@click.option("--debug_scout", is_flag=True, help="If set, generate debugging output")
@click.option(
    "--k8s", is_flag=True, help="If set, assume configuration files are annotated K8s manifests"
)
@click.option(
    "--recurse", is_flag=True, help="If set, recurse into directories below config_dir_path"
)
@click.option("--stats", is_flag=True, help="If set, dump statistics to stderr")
@click.option("--nopretty", is_flag=True, help="If set, do not pretty print the dumped JSON")
@click.option("--aconf", is_flag=True, help="If set, dump the Ambassador config")
@click.option("--ir", is_flag=True, help="If set, dump the IR")
@click.option("--xds", is_flag=True, help="If set, dump the Envoy config")
@click.option("--diag", is_flag=True, help="If set, dump the Diagnostics overview")
@click.option("--everything", is_flag=True, help="If set, dump everything")
@click.option("--features", is_flag=True, help="If set, dump the feature set")
@click.option("--profile", is_flag=True, help="If set, profile with the cProfile module")
def dump(
    config_dir_path: str,
    *,
    secret_dir_path=None,
    watt=False,
    debug=False,
    debug_scout=False,
    k8s=False,
    recurse=False,
    stats=False,
    nopretty=False,
    everything=False,
    aconf=False,
    ir=False,
    xds=False,
    diag=False,
    features=False,
    profile=False,
):
    """
    Dump various forms of an Ambassador configuration for debugging

    Use --aconf, --ir, and --envoy to control what gets dumped. If none are requested, the IR
    will be dumped.

    :param config_dir_path: Configuration directory to scan for Ambassador YAML files
    """

    if not secret_dir_path:
        secret_dir_path = "/tmp/cli-secrets"

        if not os.path.isdir(secret_dir_path):
            secret_dir_path = os.path.dirname(secret_dir_path)

    if debug:
        logger.setLevel(logging.DEBUG)

    if debug_scout:
        logging.getLogger("ambassador.scout").setLevel(logging.DEBUG)

    if everything:
        aconf = True
        ir = True
        xds = True
        diag = True
        features = True
    elif not (aconf or ir or xds or diag or features):
        aconf = True
        ir = True
        xds = True
        diag = False
        features = False

    dump_aconf = aconf
    dump_ir = ir
    dump_xds = xds
    dump_diag = diag
    dump_features = features

    od = {}
    diagconfig: Optional[EnvoyConfig] = None

    _profile: Optional[cProfile.Profile] = None
    _rc = 0

    if profile:
        _profile = cProfile.Profile()
        _profile.enable()

    try:
        total_timer = Timer("total")
        total_timer.start()

        fetch_timer = Timer("fetch resources")
        with fetch_timer:
            aconf = Config()

            fetcher = ResourceFetcher(logger, aconf)

            if watt:
                fetcher.parse_watt(open(config_dir_path, "r").read())
            else:
                fetcher.load_from_filesystem(config_dir_path, k8s=k8s, recurse=True)

        load_timer = Timer("load fetched resources")
        with load_timer:
            aconf.load_all(fetcher.sorted())

        # aconf.post_error("Error from string, boo yah")
        # aconf.post_error(RichStatus.fromError("Error from RichStatus"))

        irgen_timer = Timer("ir generation")
        with irgen_timer:
            secret_handler = NullSecretHandler(logger, config_dir_path, secret_dir_path, "0")

            ir = IR(aconf, file_checker=file_checker, secret_handler=secret_handler)

        aconf_timer = Timer("aconf")
        with aconf_timer:
            if dump_aconf:
                od["aconf"] = aconf.as_dict()

        ir_timer = Timer("ir")
        with ir_timer:
            if dump_ir:
                od["ir"] = ir.as_dict()

        xds_timer = Timer("xds")
        with xds_timer:
            if dump_xds:
                config = V3Config(ir)
                diagconfig = config
                od["xds"] = config.as_dict()
        diag_timer = Timer("diag")
        with diag_timer:
            if dump_diag:
                if not diagconfig:
                    diagconfig = V3Config(ir)
                econf = typecast(EnvoyConfig, diagconfig)
                diag = Diagnostics(ir, econf)
                od["diag"] = diag.as_dict()
                od["elements"] = econf.elements

        features_timer = Timer("features")
        with features_timer:
            if dump_features:
                od["features"] = ir.features()

        # scout = Scout()
        # scout_args = {}
        #
        # if ir and not os.environ.get("AMBASSADOR_DISABLE_FEATURES", None):
        #     scout_args["features"] = ir.features()
        #
        # result = scout.report(action="dump", mode="cli", **scout_args)
        # show_notices(result)

        dump_timer = Timer("dump JSON")

        with dump_timer:
            js = dump_json(od, pretty=not nopretty)
            jslen = len(js)

        write_timer = Timer("write JSON")
        with write_timer:
            sys.stdout.write(js)
            sys.stdout.write("\n")

        total_timer.stop()

        route_count = 0
        vhost_count = 0
        filter_chain_count = 0
        filter_count = 0
        if "xds" in od:
            for listener in od["xds"]["static_resources"]["listeners"]:
                for fc in listener["filter_chains"]:
                    filter_chain_count += 1
                    for f in fc["filters"]:
                        filter_count += 1
                        for vh in f["typed_config"]["route_config"]["virtual_hosts"]:
                            vhost_count += 1
                            route_count += len(vh["routes"])

        if stats:
            sys.stderr.write("STATS:\n")
            sys.stderr.write("  config bytes:  %d\n" % jslen)
            sys.stderr.write("  vhosts:        %d\n" % vhost_count)
            sys.stderr.write("  filter chains: %d\n" % filter_chain_count)
            sys.stderr.write("  filters:       %d\n" % filter_count)
            sys.stderr.write("  routes:        %d\n" % route_count)
            sys.stderr.write(
                "  routes/vhosts: %.3f\n" % float(float(route_count) / float(vhost_count))
            )
            sys.stderr.write("TIMERS:\n")
            sys.stderr.write("  fetch resources:  %.3fs\n" % fetch_timer.average)
            sys.stderr.write("  load resources:   %.3fs\n" % load_timer.average)
            sys.stderr.write("  ir generation:    %.3fs\n" % irgen_timer.average)
            sys.stderr.write("  aconf:            %.3fs\n" % aconf_timer.average)
            sys.stderr.write("  envoy:            %.3fs\n" % xds_timer.average)
            sys.stderr.write("  diag:             %.3fs\n" % diag_timer.average)
            sys.stderr.write("  features:         %.3fs\n" % features_timer.average)
            sys.stderr.write("  dump json:        %.3fs\n" % dump_timer.average)
            sys.stderr.write("  write json:       %.3fs\n" % write_timer.average)
            sys.stderr.write("  ----------------------\n")
            sys.stderr.write("  total: %.3fs\n" % total_timer.average)
    except Exception as e:
        handle_exception("EXCEPTION from dump", e, config_dir_path=config_dir_path)
        _rc = 1

    if _profile:
        _profile.disable()
        _profile.dump_stats("ambassador.profile")

    sys.exit(_rc)


@click.command()
@click.argument("config_dir_path", type=click.Path())
def validate(config_dir_path: str):
    """
    Validate an Ambassador configuration. This is an extension of "config" that
    redirects output to devnull and always exits on error.

    :param config_dir_path: Configuration directory to scan for Ambassador YAML files
    """
    config(config_dir_path, os.devnull, exit_on_error=True)


@click.command()
@click.argument("config_dir_path", type=click.Path())
@click.argument("output_json_path", type=click.Path())
@click.option("--debug", is_flag=True, help="If set, generate debugging output")
@click.option(
    "--debug-scout", is_flag=True, help="If set, generate debugging output when talking to Scout"
)
@click.option(
    "--check", is_flag=True, help="If set, generate configuration only if it doesn't already exist"
)
@click.option(
    "--k8s", is_flag=True, help="If set, assume configuration files are annotated K8s manifests"
)
@click.option(
    "--exit-on-error",
    is_flag=True,
    help="If set, will exit with status 1 on any configuration error",
)
@click.option(
    "--ir", type=click.Path(), help="Pathname to which to dump the IR (not dumped if not present)"
)
@click.option(
    "--aconf",
    type=click.Path(),
    help="Pathname to which to dump the aconf (not dumped if not present)",
)
def config(
    config_dir_path: str,
    output_json_path: str,
    *,
    debug=False,
    debug_scout=False,
    check=False,
    k8s=False,
    ir=None,
    aconf=None,
    exit_on_error=False,
):
    """
    Generate an Envoy configuration

    :param config_dir_path: Configuration directory to scan for Ambassador YAML files

    :param output_json_path: Path to output envoy.json
    """

    if debug:
        logger.setLevel(logging.DEBUG)

    if debug_scout:
        logging.getLogger("ambassador.scout").setLevel(logging.DEBUG)

    try:
        logger.debug("CHECK MODE  %s" % check)
        logger.debug("CONFIG DIR  %s" % config_dir_path)
        logger.debug("OUTPUT PATH %s" % output_json_path)

        dump_aconf: Optional[str] = aconf
        dump_ir: Optional[str] = ir

        # Bypass the existence check...
        output_exists = False

        if check:
            # ...oh no wait, they explicitly asked for the existence check!
            # Assume that the file exists (ie, we'll do nothing) unless we
            # determine otherwise.
            output_exists = True

            try:
                parse_json(open(output_json_path, "r").read())
            except FileNotFoundError:
                logger.debug("output file does not exist")
                output_exists = False
            except OSError:
                logger.warning("output file is not sane?")
                output_exists = False
            except json.decoder.JSONDecodeError:
                logger.warning("output file is not valid JSON")
                output_exists = False

            logger.info("Output file %s" % ("exists" if output_exists else "does not exist"))

        rc = RichStatus.fromError("impossible error")

        if not output_exists:
            # Either we didn't need to check, or the check didn't turn up
            # a valid config. Regenerate.
            logger.info("Generating new Envoy configuration...")

            aconf = Config()
            fetcher = ResourceFetcher(logger, aconf)
            fetcher.load_from_filesystem(config_dir_path, k8s=k8s)
            aconf.load_all(fetcher.sorted())

            if dump_aconf:
                with open(dump_aconf, "w") as output:
                    output.write(aconf.as_json())
                    output.write("\n")

            # If exit_on_error is set, log _errors and exit with status 1
            if exit_on_error and aconf.errors:
                raise Exception("errors in: {0}".format(", ".join(aconf.errors.keys())))

            secret_handler = NullSecretHandler(logger, config_dir_path, config_dir_path, "0")

            ir = IR(aconf, file_checker=file_checker, secret_handler=secret_handler)

            if dump_ir:
                with open(dump_ir, "w") as output:
                    output.write(ir.as_json())
                    output.write("\n")

            logger.info("Writing envoy configuration")
            config = V3Config(ir)
            rc = RichStatus.OK(msg="huh_xds")

            if rc:
                with open(output_json_path, "w") as output:
                    output.write(config.as_json())
                    output.write("\n")
            else:
                logger.error("Could not generate new Envoy configuration: %s" % rc.error)

        scout = Scout()
        result = scout.report(action="config", mode="cli")
        show_notices(result)
    except Exception as e:
        handle_exception(
            "EXCEPTION from config",
            e,
            config_dir_path=config_dir_path,
            output_json_path=output_json_path,
        )

        # This is fatal.
        sys.exit(1)


def version_callback(ctx: click.core.Context, param: click.Parameter, value: bool) -> None:
    if not value:
        return
    version()
    ctx.exit()


def showid_callback(ctx: click.core.Context, param: click.Parameter, value: bool) -> None:
    if not value:
        return
    showid()
    ctx.exit()


@click.group(
    no_args_is_help=False,
    commands=[config, dump, validate],
)
@click.option(
    "--version",
    is_flag=True,
    expose_value=False,
    callback=version_callback,
    help="Show the Emissary version number and exit.",
)
@click.option(
    "--showid",
    is_flag=True,
    expose_value=False,
    callback=showid_callback,
    help="Show the cluster ID and exit.",
)
def main():
    """Generate an Envoy config, or manage an Ambassador deployment. Use

        ambassador.py command --help

    for more help, or

        ambassador.py --version

    to see Ambassador's version.
    """
    pass


if __name__ == "__main__":
    main()
