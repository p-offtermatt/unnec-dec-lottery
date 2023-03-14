What does the app do?

It has two chains, lottery-oracle (l) and ticket-office (t).
l organizes lotteries, but tickets are bought by sending transactions to t.
The lifecycle of a round of lottery is:
* A round of lottery is started on l
* Users can enter the lottery by sending a tx on chain t and buying a ticket
* t sends ibc packets to tell l when a new user enters the lottery
* once the deadline height has passed, l will pick a winner (by a deterministic process)
and send an ibc packet with the winner to l
* l will distribute the funds to the winner, and everyone else goes empty

What do we need to store?
* on l:
  * rounds of lottery
    * ids of users that bought tickets for the round
    * the deadline (a block height)
    * an id for the round so that users can choose which round to enter
* on t:
  * rounds of lottery (as well)
    * id of the round
    * the total betting amount entered in that round

Handle this via lists:
on l:
* list of lottery rounds on l
    * deadline
    * id
    * list of users in round

on t:
* list of lottery rounds
  * id
  * total amount in the betting pool 

needed TXs:
on l:
* CreateLottery
  * deadline
  * id is auto assigned
* cancel lottery round (refund users)
  * id

on t:
* EnterLottery
  * id
  * ticket price

needed ibc packets:
from t to l:
* TicketBought (sent when a user sends EnterLottery)
  * Acknowledged if lottery round is before its deadline, or error when the round has ended (user is refunded)
l to t
* WinnerPicked (sent when a lottery deadline passes)
* CancelLottery (users are refunded)


Possible problem scenarios:
* User buys ticket, but while packet is in transit, round ends
  * No problem: user is refunded, because packet is rejected
* Cancel round while TicketBought is in transit:
  * No problem: refunds user
* Cancel round, but before tx is included deadline passes:
  * Bad luck, but cancellations can always be censored, no way around this

* Necessary steps:
  * ignite scaffold chain lottery --no-module
  * ignite scaffold chain ticket --no-module
  * cd lottery; ignite scaffold module lottery --ibc; cd ..
  * cd ticket; ignite scaffold module ticket --ibc; cd ..

ignite chain serve -c ticket.yaml

  ignite relayer configure -a \
  --source-rpc "http://0.0.0.0:26657" \
  --source-faucet "http://0.0.0.0:4500" \
  --source-port "lottery" \
  --source-version "lottery-1" \
  --source-gasprice "0.0000025stake" \
  --source-prefix "cosmos" \
  --source-gaslimit 300000 \
  --target-rpc "http://0.0.0.0:26659" \
  --target-faucet "http://0.0.0.0:4501" \
  --target-port "lottery" \
  --target-version "lottery-1" \
  --target-gasprice "0.0000025stake" \
  --target-prefix "cosmos" \
  --target-gaslimit 300000

  export FROML="--chain-id lottery --home ~/.lottery"
  export FROMT="--chain-id ticket --home ~/.ticket"