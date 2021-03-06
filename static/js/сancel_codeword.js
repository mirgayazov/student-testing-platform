function deleteCodeword() {
    document.getElementById("codeword").style.display = "none";
    document.getElementById("codeword").setAttribute('value',"");
    document.getElementById("deleteCodewordButton").style.display = "none";
    document.getElementById("addCodewordButton").style.display = "inline-block";
}

function addCodeword() {
    document.getElementById("codeword").style.display = "inline-block";
    document.getElementById("deleteCodewordButton").style.display = "inline-block";
    document.getElementById("addCodewordButton").style.display = "none";
}