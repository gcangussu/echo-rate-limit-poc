import json
import json.decoder
import sys
from dataclasses import dataclass
from datetime import datetime

from dateutil.parser.isoparser import isoparse


@dataclass
class LogEntry:
    time: datetime
    status: int
    id: str
    latency: int


log_by_id: dict[str, list[LogEntry]] = {}


def main():
    for line in sys.stdin:
        sys.stdout.write(line)

        try:
            # e.g. {"time":"2021-12-22T10:13:26.938228-03:00","status":429,"id":"B","latency":5903307}
            log_entry = json.loads(line.strip())
        except json.decoder.JSONDecodeError:
            continue

        id: str = log_entry["id"]
        log = log_by_id.get(id)
        if log is None:
            log = []
            log_by_id[id] = log

        log.append(
            LogEntry(
                time=isoparse(log_entry["time"]),
                status=log_entry["status"],
                id=id,
                latency=log_entry["latency"],
            )
        )


try:
    main()
except KeyboardInterrupt:
    pass


def analyse(log: list[LogEntry]):
    total_requests = len(log)
    total_success = len([i for i in log if i.status >= 200 and i.status < 400])
    total_time = (log[-1].time - log[0].time).total_seconds()
    rate_all = total_requests / total_time
    rate_effective = total_success / total_time
    return total_requests, total_time, rate_all, rate_effective


for id, log in log_by_id.items():
    total_requests, total_time, rate_all, rate_effective = analyse(log)
    print(f"--- {id} ---")
    print(f"Total requests: {total_requests}")
    print(f"Total time (s): {total_time}")
    print(f"All requests rate (req/s): {rate_all}")
    print(f"2xx requests rate (req/s): {rate_effective}")
    print("")
