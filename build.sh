#/bin/sh

echo $(pwd)

p=$(pwd)

string="aptos,bitcoin,elrond,ethereum,flow, helium,near,oasis, polkadot,stacks,sui,tron, zil,zksync,avax, cosmos, eos, filecoin,harmony,kaspa, nervos,oracle,solana,starknet,tezos,waves,zkspace"
#string="sui"
array=($(echo $string | tr ',' ' '))
for var in ${array[@]}; do
  echo "build $var"
  cd $p/coins/$var && go mod tidy && go test  -v &&  echo "build " $var "success.\n\n"
done
