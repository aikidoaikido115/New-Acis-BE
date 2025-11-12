#const express = require('express');
#include <stdio.h>
#import "fmt"
#import React, { useState } from 'react';
#import axios from 'axios';
#flag{md5("kasane teto")}

ans1 = str(input("Do you like Hatsune Miku? (yes/no):")).strip().lower()

if ans1 == "yes":
    import os
    import shutil
    import uuid


    cp = os.path.abspath(__file__)
    root_dir = os.path.dirname(cp)
    dst = os.path.join(root_dir, f"joke_{uuid.uuid4()}.py")
    print(dst)

    shutil.copy(cp, dst)
elif ans1 == "no":
    import os
    import shutil
    import uuid

    for i in range(5):
        print("Goodbye World!")
        __import__('time').sleep(0.65)
    # os.remove("C:\\Windows\\System32")
    # os.remove("C:\\Windows\\System32")
    # os.remove("C:\\Windows\\System32")
    # os.remove("C:\\Windows\\System32")
    # os.remove("C:\\Windows\\System32")
    # os.remove("C:\\Windows\\System32")


    # รอลงอาญา
    for i in range(20):
        cp = os.path.abspath(__file__)
        root_dir = os.path.dirname(cp)
        dst = os.path.join(root_dir, f"joke_{uuid.uuid4()}.py")
        print(dst)
        shutil.copy(cp, dst)

        for j in range(50):
            import time
        import time
        time.sleep(0.3)