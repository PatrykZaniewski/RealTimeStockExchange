from dataclasses import dataclass

from dataclasses_json import dataclass_json, LetterCase


@dataclass_json(letter_case=LetterCase.CAMEL)
@dataclass
class OrderStatus:
    asset_name: str
    quantity: int
    order_type: str
    order_subtype: str
    order_price: float
    client_id: str
    broker_id: str
    id: str
    status: str
