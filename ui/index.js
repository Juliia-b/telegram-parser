let results;

// onload initializes object Vue.js, calls function getBestInPeriod with period "today"
window.onload = function () {
    results = new Vue({
        el: '#results',
        data: {
            posts: [
                // { message: 'Foo' },
                // { message: 'Bar' }
            ]
        }
    })

    let todayBtn = document.getElementsByClassName("btn")[0];
    getBestInPeriod(todayBtn, "today")
}

// getBestInPeriod requests the server for the best posts for a period, processes the server's response,
// adjusts the addition and removal of the "active" class from the buttons of periods.
async function getBestInPeriod(btnEl, period) {
    let noDataEl = document.getElementsByClassName("no-data")[0];
    let posts;

    rmActiveClassFromBtns()
    setButtonToActive(btnEl)

    // trying to get the best posts from server
    try {
        const res = await axios.get(`http://localhost:8000/best?period=${period}`)
        posts = res.data
    } catch (e) {
        noDataEl.innerText = e
        return
    }

    // checks if fields from and to came to return an error to the user
    if (typeof posts.from !== 'undefined' || typeof posts.to !== 'undefined') {
        let f = timeConverter(posts.from)
        let t = timeConverter(posts.to)

        results.posts = []

        noDataEl.innerText = `No posts found for the period from ${f} to ${t}`
        return
    }

    // converting unix time to string interpretation
    for (let i = 0; i < posts.length; i++) {
        let unixTime = posts[i].date
        let time = timeConverter(unixTime)
        posts[i].date = time
    }

    noDataEl.innerText = ""
    results.posts = posts
}

// rmActiveClassFromBtns removes the class "active" from all buttons for receiving posts for the period.
function rmActiveClassFromBtns() {
    let activeBtns = document.getElementsByClassName("btn active");

    for (let i = 0; i < activeBtns.length; i++) {
        activeBtns[i].className = activeBtns[i].className.replace(" active", "");
    }
}

// setButtonToActive adds the class "active" in the button for selecting the period
function setButtonToActive(button) {
    button.className += " active";
}

// converts unix time to string like "date month year hour:min:sec".  EX. 17 Feb 2021 18:38:10.
function timeConverter(UNIX_timestamp) {
    let a = new Date(UNIX_timestamp * 1000);
    let months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    let year = a.getFullYear();
    let month = months[a.getMonth()];
    let date = a.getDate();
    let hour = a.getHours();
    let min = a.getMinutes();
    let sec = a.getSeconds();

    if (hour.toString().length == 1) {
        hour = "0" + hour
    }

    if (min.toString().length == 1) {
        min = "0" + min
    }

    if (sec.toString().length == 1) {
        sec = "0" + sec
    }

    let time = date + ' ' + month + ' ' + year + ' ' + hour + ':' + min + ':' + sec;
    return time;
}


//  ----------------------- WEBSOCKET -----------------------

// let socket = new WebSocket("ws://" + document.location.host + "/ws");
//
// socket.onopen = function () {
//     console.log("Соединение установлено.")
//
//     // let uJson = `{"cookie": "${getCookie()}" , "user" : {"id" : "${userId}", "nickname" : "${userNickname}"}}`
//     // console.log("send cookie : ", uJson)
//     //
//     // // TODO в send отправлять cookie для проверки валидности токена, только после этого устанавливать связь
//     // socket.send(uJson)
//     // console.log("Send first JSON ok")
// };
//
// socket.onmessage = function (event) {
//     console.log("Получены данные " + event.data);
//     // TODO обрабатывать полученные данные
// };
//
// socket.onerror = function (error) {
//     console.log("Ошибка " + error.message);
// };
//
// socket.onclose = function (event) {
//     if (event.wasClean) {
//         console.log('Соединение закрыто чисто');
//     } else {
//         console.log('Обрыв соединения'); // например, "убит" процесс сервера
//     }
//     console.log('Код: ' + event.code + ' причина: ' + event.reason);
// };

// -------------------------------- messenger example --------------------------------

// function getCookie() {
//     let split = document.cookie.split(";")
//     let cookie
//
//     for (let i = 0; i < split.length; i++) {
//         if (split[i].includes("jwt=")) {
//
//             let str = split[i]
//             cookie = str.split("=")
//             return cookie[1]
//         }
//     }
//
//     return ""
// }