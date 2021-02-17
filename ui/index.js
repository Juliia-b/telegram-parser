let results;

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

    //
    if (typeof posts.from !== 'undefined' || typeof posts.to !== 'undefined') {
        let f = timeConverter(posts.from)
        let t = timeConverter(posts.to)

        results.posts = []

        noDataEl.innerText = `No posts found for the period from ${f} to ${t}`
        return
    }

    for (let i = 0; i < posts.length; i++) {
        let unixTime = posts[i].date
        let time = timeConverter(unixTime)
        posts[i].date = time
    }

    noDataEl.innerText = ""

    results.posts = posts
}

function rmActiveClassFromBtns() {
    let activeBtns = document.getElementsByClassName("btn active");

    for (let i = 0; i < activeBtns.length; i++) {
        activeBtns[i].className = activeBtns[i].className.replace(" active", "");
    }
}

function setButtonToActive(button) {
    button.className += " active";
}

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