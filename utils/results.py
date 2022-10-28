from dataclasses import dataclass


@dataclass
class ProcessingResults:
    broker_facade_processing: int
    broker_core_db_operations: int
    broker_core_processing: int
    broker_