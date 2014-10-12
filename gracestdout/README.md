#gracestdout

命令行环境下让标准输出比较优雅方便查看..

格式:

* json
* todo


##Demo

    echo '{"text":"abc","text2":{"key1":1,"key2":"haha","key3":{"name":"obj"}},"text3":"hey"}'|gracestdout
    {
        "text": "abc",
        "text2": {
            "key1": 1,
            "key2": "haha",
            "key3": {
                "name": "obj"
            }
        },
        "text3": "hey"
    }
