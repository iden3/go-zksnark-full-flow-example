include "../node_modules/circomlib/circuits/comparators.circom";
include "../node_modules/circomlib/circuits/poseidon.circom";

template Example() {
    signal private input a;
    signal private input b;
    signal input c;
    component h = Poseidon(2, 6, 8, 57);
    h.inputs[0] <== a;
    h.inputs[1] <== b;

    component eq = IsEqual();
    eq.in[0] <== h.out;
    eq.in[1] <== c;
    eq.out === 1;
}

component main = Example();
