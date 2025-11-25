import logging

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.input import text_file

log = logging.getLogger(__name__)


class ProcCpuinfo(Source):
    @classmethod
    def __has_flags(cls, check_flags, all_flags):
        check_flags = check_flags.split()
        all_flags = all_flags.split()
        log.debug(f"count is {len(check_flags)} {len(all_flags)}")
        log.debug(f"{check_flags=}")
        log.debug(f"{all_flags=}")
        for flag in check_flags:
            if flag not in all_flags:
                return False
        return True

    def provides(self, provides: Provides) -> None:
        provides["phy.cpu.arch_version"] = self

    def requires(self, requires: Requires) -> None:
        requires.files.append("/proc/cpuinfo")

    def parse(self, facts: Facts) -> None:
        if facts["phy.arch"] not in ("x86_64", "amd64"):
            log.info("Not an x86_64/amd64 architecture, skipping cpuinfo flags parsing")
            return

        vers = [
            "lm cmov cx8 fpu fxsr mmx syscall sse2",
            "cx16 lahf_lm popcnt sse4_1 sse4_2 ssse3",
            "avx avx2 bmi1 bmi2 f16c fma abm movbe xsave",
            "avx512f avx512bw avx512cd avx512dq avx512vl",
        ]

        data = text_file("/proc/cpuinfo")
        # log.debug(f"data={data.splitlines() if data else []}") This can be HUGE
        cpu_flags = ""
        if data:
            for line in data.splitlines():
                if line.startswith("flags"):
                    cpu_flags = line.split(":", 1)[1].strip()
                    break
        log.debug(f"cpu_flags={cpu_flags.split()}")
        cpu_version = 0
        for idx, v in enumerate(vers):
            if self.__has_flags(v, cpu_flags):
                cpu_version += 1
            else:
                break
        facts["phy.cpu.arch_version"] = f"x86_64_v{cpu_version}"
