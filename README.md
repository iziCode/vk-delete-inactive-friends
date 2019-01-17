### Для того чтобы удалить неактивных друзей нужно:
#### 1) Получить access_token, открыв для этого ссылку:
https://oauth.vk.com/authorize?client_id=6819511&scope=friends&redirect_uri=https://api.vk.com/blank.html&display=page&response_type=token&v=5.92
##### приложение попросит разрешение, нажмите разрешить, затем скопируйте с поисковой строки значение access_token(), пример строки:
https://api.vk.com/blank.html#access_token=60d5db9be16df4ea3b1e4241b246942832bbff5ef33b0127345ef7592b0dc957b64c78d22eae18ce1542&expires_in=86400&user_id=145465502
##### Вам нужно это значение:
60d5db9be16df4ea3b1e4241b246942832bbff5ef33b0127345ef7592b0dc957b64c78d22eae18ce1542
#### 2) Запустить программу  «vk-delete-inactive-friends.exe»
#### 3) Ввести access_token и кол-во недель которое не заходил пользователь в сеть 
#### (Например, вы ввели 7, следовательно все друзья которые не заходили 7 недель и более будут удалены)
