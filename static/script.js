let verified = false;
let allCorrect = 0;
let noCorrect = 0;
let noIncorrect = 0;

function activate(answer) {
    if(verified) {
        return;
    }

    if(answer.className.includes("active")) {
        answer.className = answer.className.replaceAll(" active", "");
    } else {
        answer.className = answer.className + " active";
    }
}

function verify() {
    if(!verified) {
        verified = true;
    } else {
        return;
    }

    document.getElementById("btn-verify").setAttribute("disabled", "true");

    var answers = document.getElementsByClassName("answer")
    for(let i = 0; i < answers.length; i++) {
        var ans = answers[i];

        if(ans.className.includes(" correct")) {
            allCorrect++;
        }

        if(ans.className.includes(" correct") && ans.className.includes("active")) {
            ans.className = ans.className + " bg-success";
            noCorrect++;
        }

        if(ans.className.includes(" correct") && !ans.className.includes("active")) {
            ans.className = ans.className + " bg-warning";
        }

        if(ans.className.includes("incorrect") && ans.className.includes("active")) {
            ans.className = ans.className + " bg-danger";
            noIncorrect++;
        }

        ans.className = ans.className.replaceAll("active", " ");
    }

    var scoreAlert = document.getElementById("score-alert")

    scoreAlert.className = scoreAlert.className.replaceAll("alert-primary");
    if(noCorrect == allCorrect && noIncorrect == 0) {
        scoreAlert.className = scoreAlert.className + " alert-success";
        scoreAlert.innerHTML = "Congrats, all anwers are correct!";
    } else if(noCorrect > 0) {
        scoreAlert.className = scoreAlert.className + " alert-warning";
        scoreAlert.innerHTML = "You are almost correct.";
    } else {
        scoreAlert.className = scoreAlert.className + " alert-danger";
        scoreAlert.innerHTML = "None of your answer was correct.";

    }

    var verifyBtn = document.getElementById("verify-btns");
    verifyBtn.className = verifyBtn.className + " d-none";

    var nextQuestionBtn = document.getElementById("next-question-btns");
    nextQuestionBtn.className = nextQuestionBtn.className.replaceAll("d-none");
}