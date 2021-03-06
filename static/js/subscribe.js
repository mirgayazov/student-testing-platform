function openForm(id, courseName) {
    document.getElementById("myForm").style.display = "block";
    document.getElementById("courseID").setAttribute('value',id)
    document.getElementById("label").innerText="Подписка на курс: "+courseName
}

function closeForm() {
    document.getElementById("myForm").style.display = "none";
}