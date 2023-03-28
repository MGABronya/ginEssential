# ginEssential

- ## 用户相关

  - **接口地址：/regist**

    **功能：用于用户注册**

    **方法类型：POST**

    请求参数：Body部分，form-data类型，接收四个字符串分别为Email，Password，Name，Verify。其中Email需要符合邮箱格式，Password最少需要六位，Name最多为20位长度，Verify必须与邮箱验证码相同，注意，用户的邮箱和名称均不能重复。

    可能的返回值：成功则返回200与token，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/verify/:id**

    **功能：用于用户请求验证邮箱**

    **方法类型：GET**

    请求参数：需要在接口地址部分（:id）给出用户邮箱

    返回值：成功则返回200并向相应邮箱发送验证邮件，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/security**

    **功能：用于用户找回密码**

    **方法类型：PUT**

    请求参数：Body部分，form-data类型，接收两个字符串分别为Email，Verify。其中Verify必须与邮箱验证码相同。

    返回值：成功则返回200并向相应邮箱发送重置后的密码，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/updatepass**

    **功能：用于用户修改密码**

    **方法类型：PUT**

    请求参数：Body部分，form-data类型，接收两个字符串分别为first，second。其中first为旧密码，second为新密码。

    返回值：成功则返回200并修改数据库中的密码，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/login**

    **功能：用于用户登录**

    **方法类型：POST**

    请求参数：Body部分，form-data类型，接收两个字符串分别为Email，Password。其中Email需要符合邮箱格式，Password最少需要六位。

    可能的返回值：成功则返回200与token，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/personal**

    **功能：获取当前用户的个人信息**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：返回用户的一些个人信息，格式为json包含Name,Email,Telephone,Blog,QQ,Sex,Address,Hobby,Icon,BackGround，除Sex为bool类型外，其它均为字符串类型，其中Icon和BackGround表示头像和背景图片的文件名。

  - **接口地址：/personal/articles**

    **功能：获取当前用户的文章列表**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：返回中包含一个articles数组，其中包含了该用户所有文章的信息。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

  - **接口地址：/personal/posts**

    **功能：获取当前用户的帖子列表**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇帖子，默认值为20）。

    返回值：返回中包含一个posts数组，其中包含了该用户所有帖子的信息。

  - **接口地址：/personal/threads**

    **功能：获取当前用户的跟帖列表**

    **方法：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇跟帖，默认值为20）。

    返回值：返回中包含一个threads数组，其中包含了该用户所有跟帖的信息。

  - **接口地址：/personal/users**

    **功能：提供一组userid，返回一组对应的用户信息**

    **方法：POST**

    请求参数：在Body，raw格式给出json类型数据包含一个userId，userId为一个数组，数组中包含了需要查询的user组。

    返回值：返回一个users数组，其中的用户信息与userId数组一一对应。

  - **接口地址：/personal**

    **功能：用户修改个人信息**

    **方法类型：PUT**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body	中，raw格式提供json包含Newname,Newtelephone,Newemail,Newhobby,Newaddress,Newsex,Newqq,Newblog，Verify表示修改后的信息，当email被修改时，需要提供Verify，其它情况不需要提供。

    返回值：更新成功后返回用户更新后的个人信息，否则返回错误原因。

  - **接口地址：/personal/icon**

    **功能：修改个人头像**

    **方法类型：PUT**

    请求参数：.Header中需要包含Content-Type，指名为multipart/form-data。在Body中给用form-data格式给出file（文件类型）。uthorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：更新成功后返回用户更新后的头像文件名，否则返回错误原因。

  - **接口地址：/background/show**

    **功能：用于获取用户自己使用的背景**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：用户使用的背景图片的文件名

  - **接口地址：/background/update/:id**

    **功能：用户选择其它的默认背景**

    **方法类型：PUT**

    请求参数：用户的token以及用户选择的背景的文件名（即接口地址id部分）

    返回值：更新成功时返回成功信息和文件名

  - **接口地址：/background/create**

    **功能：用户上传自定义的背景**

    **方法类型：PUT**

    请求参数：.Header中需要包含Content-Type，指名为multipart/form-data。uthorization中的Bearer Token中提供注册、登录时给出的token，同时在Body部分，form-data格式，接收file（文件类型）

    返回值：更新成功时返回成功信息和文件名

  - **接口地址：/personal/:id**

    **功能：查看某一用户的信息**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。被查看用户的id（接口地址中的id）

    返回值：成功时返回被查看用户的信息，失败时给出失败原因。
    
  - **接口地址：/personal/articles/:id**

    **功能：展示指定用户的文章列表**

    **方法：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。指定用户的用户id（即接口地址id部分），在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇帖子，默认值为20）。

    返回值：返回一个articles表示指定用户的文章列表

  - **接口地址：/personal/posts/:id**

    **功能：展示指定用户的帖子列表**

    **方法：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。指定用户的用户id（即接口地址id部分），在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇帖子，默认值为20）。

    返回值：返回一个posts表示指定用户的帖子列表

  - **接口地址：/personal/threads/:id**

    **功能：展示指定用户的跟帖列表**

    **方法：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。指定用户的用户id（即接口地址id部分），在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇帖子，默认值为20）。

    返回值：返回一个posts表示指定用户的跟帖列表

- ## 文章相关

  - **接口地址：/article**

    **功能：文章发布**

    **方法类型：POST**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选），其中title表示文章标题，content表示文章内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功时返回创建成功相关信息和文章信息article，否则给出失败原因

  - **接口地址：/article/:id**

    **功能：文章查看**

    **方法类型：GET**

    请求参数：文章的uuid（在接口地址的id处），Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功找到文章时，将会以json格式给出文章article，article中包含id,user_id,content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/article/:id**

    **功能：文章更新**

    **方法类型：PUT**

    请求参数：文章的uuid（在接口地址的id处）。uthorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选），其中title表示文章标题，content表示文章内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：成功更新文章时，将会以json格式给出文章article和，article中包含id,user_id,content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/article/:id**

    **功能：文章删除**

    **方法类型：DELETE**

    请求参数：文章的uuid（在接口地址的id处）。uthorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：成功删除文章时，将会以json格式给出文章article，article中包含id,user_id,content,create_at,updated_at,res_short,res_long。如果失败则返回失败原因。

  - **接口地址：/article/pagelist**

    **功能：给出文章列表**

    **方法类型：GET**

    请求参数：uthorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇文章，默认值为20）。

    返回值：成功时，以json格式返回两个数组articles，articles返回了相应列表的文章信息（按照创建时间排序，越早创建排序越前），如果失败则返回失败原因。

- ## 帖子相关

  - **接口地址：/post**

    **功能：创建帖子**

    **方法类型：POST**

    请求参数：uthorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含title、content、res_long(可选)、res_short（可选），其中title表示帖子标题，content表示帖子内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：创建成功时返回成功信息和帖子信息post，失败时返回失败原因

  - **接口地址：/post/:id**

    **功能：更新帖子**

    **方法类型：PUT**

    请求参数：给出帖子的id(在接口地址的:id中)。uthorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示跟帖内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：更新成功时返回成功信息，失败时返回失败原因

  - **接口地址：/post/:id**

    **功能：查看帖子**

    **方法类型：GET**

    请求参数：thorization中的Bearer Token中提供注册、登录时给出的token。在接口地址的id中给出帖子的id。

    返回值：成功时返回帖子post。post中包含了id,user_id,title,content,res_long,res_short,create_at,updated_at。

  - **接口地址：/post/:id**

    **功能：删除帖子**

    **方法类型：DELETE**

    请求参数：thorization中的Bearer Token中提供注册、登录时给出的token。在接口地址的id中给出帖子的id。

    返回值：成功时返回帖子post。失败时返回错误信息。

  - **接口地址：/post/pagelist**

    **功能：提供帖子列表**

    **方法类型：GET**

    请求参数：uthorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇帖子，默认值为20）。

    返回值：成功时，以json格式返回两个数组posts，posts返回了相应列表的帖子信息（按照创建时间排序，越早创建排序越前）

- ## 跟帖相关

  - **接口地址：/thread/:id**

    **功能：创建跟帖**

    **方法类型：POST**

    请求参数：给出主贴的id(在接口地址的:id中)。uthorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示跟帖内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：创建成功时返回成功信息，失败时返回失败原因

  - **接口地址：/thread/:id**

    **功能：更新跟帖**

    **方法类型：PUT**

    请求参数：给出跟帖的id(在接口地址的:id中)。uthorization中的Bearer Token中提供注册、登录时给出的token。在Body，raw格式给出json类型数据包含content、res_long(可选)、res_short（可选），其中content表示跟帖内容，res_long表示长文本备用键值，res_short表示短文本备用键值。

    返回值：更新成功时返回成功信息，失败时返回失败原因

    **功能：查看跟帖**

    **方法类型：GET**

    请求参数：thorization中的Bearer Token中提供注册、登录时给出的token。在接口地址的id中给出帖子的id。

    返回值：成功时返回帖子thread，thread中包含了id,user_id,post_id,content,res_long,res_short,create_at,updated_at

  - **接口地址：/thread/:id**

    **功能：删除跟帖**

    **方法类型：DELETE**

    请求参数：thorization中的Bearer Token中提供注册、登录时给出的token。在接口地址的id中给出跟帖的id。

    返回值：成功时返回跟帖thread。失败时返回错误信息。

  - **接口地址：/thread/pagelist/:id**

    **功能：提供跟帖列表**

    **方法类型：GET**

    请求参数：uthorization中的Bearer Token中提供注册、登录时给出的token。在Params处提供pageNum（表示第几页，默认值为1）和pageSize（表示一页多少篇跟帖，默认值为20）,在接口id处给出帖子的id。

    返回值：成功时，以json格式返回一个数组threads，threads返回了相应列表的跟帖信息（按照创建时间排序，越早创建排序越前）

