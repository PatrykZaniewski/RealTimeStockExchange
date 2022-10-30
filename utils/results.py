from dataclasses import dataclass


@dataclass
class ProcessingResults:
    client_broker_facade_communication: int
    broker_facade_processing: int
    broker_facade_broker_core_communication: int
    broker_core_order_db_operations: int
    broker_core_order_processing: int
    broker_core_broker_order_executor_communication: int
    broker_order_executor_processing: int
    broker_order_executor_stock_order_collector_communication: int
    stock_order_collector_processing: int
    stock_order_collector_stock_core_communication: int
    stock_core_db_operations: int
    stock_core_processing: int
    stock_core_broker_order_status_collector_communication: int
    broker_order_status_collector_processing: int
    broker_order_status_collector_broker_core_communication: int
    broker_core_status_db_operations: int
    broker_core_status_processing: int
    broker_core_broker_data_streamer_communication: int
    broker_data_streamer_processing: int
    broker_data_streamer_client_communication: int
