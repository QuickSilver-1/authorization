# authorization
  
  ![Build Status](https://github.com/QuickSilver-1/authorization/actions/workflows/go.yml/badge.svg)

 <h3>Прототип авторизационного сервера</h3>

<h2>Запуск</h2>

<h3>Запуск сервера</h3>
<ul>
<li>Сначала необходимо настроить <code>cmd/auth/config.env</code> для запуска сервера на своих данных (сейчас он настроен на работу с моим сервером, можно использовать его для тестов)</li> 
</ul>
  <h4>После можно запустить сервер одним из 3 способов</h4>
<ol>
<li>Если на компьютере установлена утилита <code>Make</code>, то можно использовать Makefile для запуска <code>make up</code></li>
<li>Если <code>Make</code> не установлена, то необоходимо вручную собрать и запустить докер образ <code>docker build . --file Dockerfile -t app:latest</code>---><code>docker-compose up</code></li>
<li>Также можно запустить код без докера <code>cd cmd/auth</code>---><code>go run .</code></li>
</ol>

<h2>Общее описание</h2>
Спасибо за необычную задачу, было интересно её делать. Реализовал все необходимые функции из ТЗ. Используемая база данных - postgresql. <h3>Спецификацию Rest api можно посмотреть в папке <code>api</code>
