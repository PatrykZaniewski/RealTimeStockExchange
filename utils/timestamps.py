from dataclasses import dataclass


@dataclass
class Timestamps:
    client_order_send: int
    broker_facade_order_received: int
    broker_facade_order_sending: int
    broker_facade_order_send: int
    broker_core_order_received: int
    broker_core_order_processing: int
    broker_core_order_processed: int
    broker_core_order_sending: int
    broker_core_order_send: int
    broker_order_executor_order_received: int
    broker_order_executor_order_sending: int
    broker_order_executor_order_send: int
    stock_order_collector_order_received: int
    stock_order_collector_order_sending: int
    stock_core_order_received: int
    stock_core_order_processing: int
    stock_core_order_processed: int
    stock_core_order_status_sending: int
    stock_core_order_status_send: int
    broker_order_status_collector_status_received: int
    broker_order_status_collector_status_sending: int
    broker_order_status_collector_status_send: int
    broker_core_status_received: int
    broker_core_status_processing: int
    broker_core_status_processed: int
    broker_core_status_sending: int
    broker_data_streamer_status_received: int
    broker_data_streamer_status_sending: int
    broker_data_streamer_status_send: int
    client_status_received: int