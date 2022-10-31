import json
from typing import List, Dict

from utils.results import ProcessingResults
from utils.timestamps import ProcessingTimestamps


def read_specified_cloud_timestamp(filename, processing_timestamps: Dict[str, ProcessingTimestamps]):
    f = open(filename)
    data = json.load(f)
    for obj in data:
        line = obj['textPayload']
        log = line.split()[2]
        log_data = log.split(",")
        order_id = log_data[0]
        if processing_timestamps.get(order_id) is None:
            processing_timestamps[order_id] = ProcessingTimestamps()
        related_processing_timestamp = processing_timestamps[order_id]
        timestamp = log_data[3]
        field_name = f'{log_data[1].lower()}_{log_data[2].lower()}'
        setattr(related_processing_timestamp, field_name, int(timestamp))


def read_client_timestamp(filename, processing_timestamps: Dict[str, ProcessingTimestamps]):
    f = open(filename)
    for log_line in f:
        log_data = log_line.split()[2].split(",")
        order_id = log_data[0]
        if processing_timestamps.get(order_id) is None:
            processing_timestamps[order_id] = ProcessingTimestamps()
        related_processing_timestamp = processing_timestamps[order_id]
        timestamp = log_data[3]
        field_name = f'{log_data[1].lower()}_{log_data[2].lower()}'
        setattr(related_processing_timestamp, field_name, int(timestamp))


def read_timestamps():
    processing_timestamps: Dict[str, ProcessingTimestamps] = {}
    read_specified_cloud_timestamp('broker.json', processing_timestamps)
    read_specified_cloud_timestamp('stock.json', processing_timestamps)
    read_client_timestamp('client.json', processing_timestamps)
    return processing_timestamps


def generate_processing_results(timestamps: List[ProcessingTimestamps]):
    results = []
    for timestamp in timestamps:
        results.append(ProcessingResults(
            client_broker_facade_communication=int(
                timestamp.broker_facade_order_received - timestamp.client_order_send),
            broker_facade_processing=int(((
                                                      timestamp.broker_facade_order_send + timestamp.broker_facade_order_sending) / 2) - timestamp.broker_facade_order_received),
            broker_facade_broker_core_communication=int(timestamp.broker_core_order_received - (
                        (timestamp.broker_facade_order_send + timestamp.broker_facade_order_sending) / 2)),
            broker_core_order_db_operations=int(
                timestamp.broker_core_order_processed - timestamp.broker_core_order_processing),
            broker_core_order_processing=int((
                                                         timestamp.broker_core_order_send + timestamp.broker_core_order_sending) / 2 - timestamp.broker_core_order_received),
            broker_core_broker_order_executor_communication=int(timestamp.broker_order_executor_order_received - (
                        (timestamp.broker_core_order_sending + timestamp.broker_core_order_send) / 2)),
            broker_order_executor_processing=int(((
                                                              timestamp.broker_order_executor_order_sending + timestamp.broker_order_executor_order_send) / 2) - timestamp.broker_order_executor_order_received),
            broker_order_executor_stock_order_collector_communication=int(
                timestamp.stock_order_collector_order_received - ((
                                                                              timestamp.broker_order_executor_order_sending + timestamp.broker_order_executor_order_send) / 2)),
            stock_order_collector_processing=int((
                                                             timestamp.stock_order_collector_order_sending + timestamp.stock_order_collector_order_send) / 2 - timestamp.stock_order_collector_order_received),
            stock_order_collector_stock_core_communication=int(timestamp.stock_core_order_received - ((
                                                                                                                  timestamp.stock_order_collector_order_sending + timestamp.stock_order_collector_order_send) / 2)),
            stock_core_db_operations=int(timestamp.stock_core_order_processed - timestamp.stock_core_order_processing),
            stock_core_processing=int(((
                                                   timestamp.stock_core_status_sending + timestamp.stock_core_status_send) / 2) - timestamp.stock_core_order_received),
            stock_core_broker_order_status_collector_communication=int(
                timestamp.broker_order_status_collector_status_received - (
                            (timestamp.stock_core_status_sending + timestamp.stock_core_status_send) / 2)),
            broker_order_status_collector_processing=int(((
                                                                      timestamp.broker_order_status_collector_status_sending + timestamp.broker_order_status_collector_status_send) / 2) - timestamp.broker_order_status_collector_status_received),
            broker_order_status_collector_broker_core_communication=int(timestamp.broker_core_status_received - ((
                                                                                                                             timestamp.broker_order_status_collector_status_sending + timestamp.broker_order_status_collector_status_send) / 2)),
            broker_core_status_db_operations=int(
                timestamp.broker_core_status_processed - timestamp.broker_core_status_processing),
            broker_core_status_processing=int(((
                                                           timestamp.broker_core_status_sending + timestamp.broker_core_status_send) / 2) - timestamp.broker_core_status_received),
            broker_core_broker_data_streamer_communication=int(timestamp.broker_data_streamer_status_received - (
                        (timestamp.broker_core_status_sending + timestamp.broker_core_status_send) / 2)),
            broker_data_streamer_processing=int(((
                                                             timestamp.broker_data_streamer_status_sending + timestamp.broker_data_streamer_status_send) / 2) - timestamp.broker_data_streamer_status_received),
            broker_data_streamer_client_communication=int(timestamp.client_status_received - ((
                                                                                                          timestamp.broker_data_streamer_status_sending + timestamp.broker_data_streamer_status_send) / 2))
        ))
    return results


def run_calculations():
    timestamps: Dict[str, ProcessingTimestamps] = read_timestamps()
    results = generate_processing_results(list(timestamps.values()))


if __name__ == "__main__":
    run_calculations()
