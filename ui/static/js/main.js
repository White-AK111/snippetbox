"use strict";

//Добавляем класс live элементам навигации
let navLinks = document.querySelectorAll("nav a");
for (let i = 0; i < navLinks.length; i++) {
    let link = navLinks[i]
    if (link.getAttribute('href') == window.location.pathname) {
        link.classList.add("live");
        break;
    }
}

//Форматируем дату в человеческий вид
let utcDates = document.querySelectorAll(".utcDate");
for (let i = 0; i < utcDates.length; i++) {
    let cDate = new Date(Date.parse(rectifyFormat(utcDates[i].innerText)));
    utcDates[i].innerText = formatDate(cDate);
}

//Форматирование даты к виду dd.MM.yyyy hh:mm
function formatDate(date) {
    let dd = date.getDate();
    if (dd < 10) dd = '0' + dd;

    let MM = date.getMonth() + 1;
    if (MM < 10) MM = '0' + MM;

    let yyyy = date.getFullYear();

    let hh = date.getHours();
    if (hh < 10) hh = '0' + hh;

    let mm = date.getMinutes();
    if (mm < 10) mm = '0' + mm;

    return dd + '.' + MM + '.' + yyyy + " " + hh + ":" + mm;
}

//Format date from pgtype.Timestamptz
function rectifyFormat(s) {
    let b = s.split(" ");
    return b[0].substr(1,10) + 'T'  + b[1].substr(0,12) + b[2];
}