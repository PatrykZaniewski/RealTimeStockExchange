from dataclasses import dataclass
from dataclasses_json import dataclass_json

from utils.results import ProcessingResults


@dataclass_json
@dataclass
class AggregatedProcessingResults:
    min_time: ProcessingResults = ProcessingResults()
    max_time: ProcessingResults = ProcessingResults()
    avg_time: ProcessingResults = ProcessingResults()
    median_time: ProcessingResults = ProcessingResults()
