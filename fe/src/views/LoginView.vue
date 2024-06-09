<template>
    <body id="poster">
      <h1 style="text-align: center;">管理员登录</h1>

        <el-form class="login-container" label-position="left" label-width="0px">

            <h3 class="login_title">
                系统登录
                <el-button @click="toRegister()">点我注册</el-button>
            </h3>
            
            <el-form-item label="">
                <el-input type="text" v-model="loginForm.loginName" placeholder="用户名" autocomplete="off"></el-input>
            </el-form-item>

            <el-form-item label="">
                <el-input type="password" v-model="loginForm.password" placeholder="密码" autocomplete="off"></el-input>
            </el-form-item>

            <el-form-item style="text-align: center;width: 100%;">
                <el-button type="primary" style="width:100%;background-color: #505458;border: none;" v-on:click="Login()">登录</el-button>
            </el-form-item>

        </el-form>
    </body>
</template>

<script>
import axios from 'axios'

  export default {
    name: 'LoginView',

    data() {
      return {
        loginForm: {
          loginName: '',
          password: '',
        },
      }
    },

    methods: {
      async Login() {
        try {
          if (
            this.loginForm.loginName &&
            this.loginForm.password
          ){
            //调用函数将数据发送到后端
          await this.sendDataToBackend();

          this.$message({
            message: '成功登录，欢迎来到用户管理系统',
            type: 'success'
          });

          this.$router.push({
            path:'/Home',
            query: {
              loginName: this.loginForm.loginName,
            }
          });
          } else {
            this.$message.error('请填写完整的登录信息');
          }
        } catch (error) {
          this.$message.error('登陆失败，请检查用户名或密码是否正确');
          console.error(error);
        } 
      },

      toRegister() {
        this.$router.push({path:'/RegisterView'})
      },

      async sendDataToBackend() {
        try {
          const response = await axios.post('http://admin.abtxw.com/api/verifyLoginData',this.loginForm,{
            headers: {
              'Content-Type': 'application/json',
            }
          });

          //处理后端响应
          console.log(response.data);
        } catch (error) {
          //处理错误
          console.error(error);
          //将错误抛出
          throw error;
        }
      }
    }
  }
</script>

<style>
    #poster {
        background-position: center;
        height: 100%;
        width: 100%;
        background-size: cover;
        position: fixed;
    }

    body {
        margin: 0px;
        padding: 0px;
    }

    .login-container {
        border-radius: 15px;
        background-clip: padding-box;
        margin: 90px auto;
        width:350px;
        padding: 35px 35px 15px 35px;
        background: #fff;
        border:1px solid #eaeaee;
        box-shadow: 0 0 25px #cac6c6;
    }
    .login_title {
        margin: 0px auto 40px auto;
        text-align: center;
    }
</style>