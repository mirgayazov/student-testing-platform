function ShowStudents() {
    document.getElementById("Users").style.display = "block";
    document.getElementById("dontShowButton").style.display = "inline-block";
    document.getElementById("ShowButton").style.display = "none";
}

function DontShowStudents() {
    document.getElementById("Users").style.display = "none";
    document.getElementById("dontShowButton").style.display = "none";
    document.getElementById("ShowButton").style.display = "inline-block";
}

function ShowData() {
    document.getElementById("Data").style.display = "block";
    document.getElementById("dontShowDataButton").style.display = "inline-block";
    document.getElementById("ShowDataButton").style.display = "none";
}

function DontShowData() {
    document.getElementById("Data").style.display = "none";
    document.getElementById("dontShowDataButton").style.display = "none";
    document.getElementById("ShowDataButton").style.display = "inline-block";
}

function DeleteQuestion(id) {
    document.getElementById("DatItem"+id).style.display = "none";
}