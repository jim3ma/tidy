use tidy;

function heredoc(fn) {
    return fn.toString().split('\n').slice(1,-1).join('\n') + '\n'
}

version = 0;

// Welcome mail
var tmpl = heredoc(function(){/*
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=2.0">
</head>
<body>
    <h2>{{ .UserName }}欢迎来到Tidy</h2>
    <br/>
</body>
</html>
*/});

db.sys_mail_template.insert({
    type: 0,
    subject: "{{ .UserName }}, 欢迎来到Tidy",
    content: tmpl,
    version: version,
});

// Reset password mail
var tmpl = heredoc(function(){/*
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=2.0">
</head>
<body>
    Hi, {{ .User.UserName }}
    <br />
    <a href="http://tf.ctidy.com/user/password.html?auth_token={{- .AuthToken -}}">点击这里重置密码</a>
    <br />
</body>
</html>
*/});

db.sys_mail_template.insert({
    type: 1,
    subject: "重置登录密码",
    content: tmpl,
    version: version,
});
