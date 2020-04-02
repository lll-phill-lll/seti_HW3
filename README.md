ИГРА шахматы

Протокол:

Registration:
<pre>
CLIENT                                                                      SERVER

REG please LOGIN %login% PASSW  %password%    ---->    Регистрация пользователя, если логин уникальный
                                              <----    SUCC REG, если регистрация успешна / ERROR NON_UNIQUE_LOGIN если логин уже есть 

</pre>

Login:
<pre>
CLIENT                                                                      SERVER

LOGIN %login% PASSW %password%     ---->    аутентификация пользователя, 
                                   <----    SUCC AUTH, если аутентификация прошла успешно /
                                            один из ответов в зависимости от ситуации: /
                                            ERROR NO_LOGIN, ERROR NO_PASSWORD, ERROR LOGIN_OR_PASSWORD"

</pre>

Создание комнаты (возможно только после регистрации или аутентификации, иначе вернется ERROR UNKNOWN_COMMAND):
<pre>
CLIENT                                                                      SERVER

LOGIN %login% PASSW %password%     ---->    аутентификация пользователя, 
                                   <----    SUCC AUTH, если аутентификация прошла успешно /
                                            один из ответов в зависимости от ситуации: /
                                            ERROR NO_LOGIN, ERROR NO_PASSWORD, ERROR LOGIN_OR_PASSWORD"

</pre>

