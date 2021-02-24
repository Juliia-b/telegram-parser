let results;
let latestTop;
let errorEl;

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

    latestTop = new Vue({
        el: '#topIn3Hours',
        data: {
            posts: [
                // { message: 'Foo' },
                // { message: 'Bar' }
            ]
        }
    })

    let todayBtn = document.getElementsByClassName("btn")[0];
    errorEl = document.getElementsByClassName("error")[0];
    getBestInPeriod(todayBtn, "today")
    getTopIn3Hours()
    startWebsocket()
}

/*-----------------------------------HANDLERS--------------------------------------*/

// getBestInPeriod requests the server for the best posts for a period, processes the server's response,
// adjusts the addition and removal of the "active" class from the buttons of periods.
async function getBestInPeriod(btnEl, period) {
    let posts;

    rmActiveFromBtns()
    setActiveToBtn(btnEl)

    // trying to get the best posts from server
    try {
        const res = await axios.get(`http://localhost:8000/best?period=${period}`)
        posts = res.data
    } catch (e) {
        errorEl.innerText = e
        return
    }

    // checks if fields from and to came to return an error to the user
    if (typeof posts.from !== 'undefined' || typeof posts.to !== 'undefined') {
        let f = timeConverter(posts.from)
        let t = timeConverter(posts.to)

        results.posts = []

        errorEl.innerText = `No posts found for the period from ${f} to ${t}`
        return
    }

    // converting unix time to string interpretation
    for (let i = 0; i < posts.length; i++) {
        let unixTime = posts[i].date
        let time = timeConverter(unixTime)
        posts[i].date = time
    }

    errorEl.innerText = ""
    results.posts = posts
}

// getTopIn3Hours returns the best posts in the last 3 hours.
// The field is self-updating by the websocket connection.
async function getTopIn3Hours() {
    let posts;

    try {
        const res = await axios.get(`http://localhost:8000/best/3hour`)
        posts = res.data
    } catch (e) {
        console.log("error when get bestIn3hour with err : ", e)
        errorEl.innerText = e
        return
    }

    updateTopIn3Hours(posts)
}

/*-----------------------------------HELPERS---------------------------------------*/

// updateTopIn3Hours replaces the content of the view block 'topIn3Hours'.
function updateTopIn3Hours(posts) {
    latestTop.posts = []
    latestTop.posts = posts
}

// changeArrow changes the position of the arrow (up / down).
function changeArrow(arrow) {
    arrow.innerHTML === '<i class="fas fa-caret-down"></i>' ? arrow.innerHTML = '<i class="fas fa-caret-up"></i>' : arrow.innerHTML = '<i class="fas fa-caret-down"></i>'
}

// openCloseText adds and removes property '-webkit-line-clamp' from element with class 'with-max-line-count'.
function openCloseText(element) {
    var textEl = element.parentNode.parentNode.children[1]
    textEl.style["webkitLineClamp"] === "3" ? textEl.style["webkitLineClamp"] = "1000" : textEl.style["webkitLineClamp"] = "3"
}

// openCloseTopText removes or adds the limit on the number of lines in an element.
function openCloseTopText(element) {
    element.style["webkitLineClamp"] === "3" ? element.style["webkitLineClamp"] = "1000" : element.style["webkitLineClamp"] = "3"
}

// rmActiveFromBtns removes the class "active" from all buttons for receiving posts for the period.
function rmActiveFromBtns() {
    let activeBtns = document.getElementsByClassName("btn active");

    for (let i = 0; i < activeBtns.length; i++) {
        activeBtns[i].className = activeBtns[i].className.replace(" active", "");
    }
}

// setActiveToBtn adds the class "active" in the button for selecting the period
function setActiveToBtn(button) {
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

    if (hour.toString().length === 1) {
        hour = "0" + hour
    }

    if (min.toString().length === 1) {
        min = "0" + min
    }

    if (sec.toString().length === 1) {
        sec = "0" + sec
    }

    return date + ' ' + month + ' ' + year + ' ' + hour + ':' + min + ':' + sec;
}

/*----------------------------------WEBSOCKET---------------------------------------*/

function startWebsocket() {
    var ws = new WebSocket("ws://" + document.location.host + "/ws")

    ws.onopen = function () {
        console.log("Соединение установлено.")
    };

    ws.onmessage = function (event) {
        console.log("Получены данные " + event.data);

        updateTopIn3Hours(event.data)
    };

    ws.onerror = function (error) {
        console.log("Ошибка " + error.message);
    };
    ws.onclose = function () {
        // connection closed, discard old websocket and create a new one in 5s
        ws = null
        setTimeout(startWebsocket, 5000)
    }
}