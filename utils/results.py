from dataclasses import dataclass, field
from dataclasses_json import dataclass_json
from typing import Optional


@dataclass_json
@dataclass
class ProcessingResults:
    client_broker_facade_communication: Optional[int] = field(default=None)
    broker_facade_processing: Optional[int] = field(default=None)
    broker_facade_broker_core_communication: Optional[int] = field(default=None)
    broker_core_order_db_operations: Optional[int] = field(default=None)
    broker_core_order_processing: Optional[int] = field(default=None)
    broker_core_broker_order_executor_communication: Optional[int] = field(default=None)
    broker_order_executor_processing: Optional[int] = field(default=None)
    broker_order_executor_stock_order_collector_communication: Optional[int] = field(default=None)
    stock_order_collector_processing: Optional[int] = field(default=None)
    stock_order_collector_stock_core_communication: Optional[int] = field(default=None)
    stock_core_db_operations: Optional[int] = field(default=None)
    stock_core_processing: Optional[int] = field(default=None)
    stock_core_broker_order_status_collector_communication: Optional[int] = field(default=None)
    broker_order_status_collector_processing: Optional[int] = field(default=None)
    broker_order_status_collector_broker_core_communication: Optional[int] = field(default=None)
    broker_core_status_db_operations: Optional[int] = field(default=None)
    broker_core_status_processing: Optional[int] = field(default=None)
    broker_core_broker_data_streamer_communication: Optional[int] = field(default=None)
    broker_data_streamer_processing: Optional[int] = field(default=None)
    broker_data_streamer_client_communication: Optional[int] = field(default=None)
