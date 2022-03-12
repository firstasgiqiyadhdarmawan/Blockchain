pragma soldity ^0.5.0

contract helloWorld{

    uint256 totalCoin:

    function addCoin(uint256 coin) public{
        totalCoin += coin;
    }

    function getCoin() public view returns(uint256){
        return totalCoin;
    }
    
}