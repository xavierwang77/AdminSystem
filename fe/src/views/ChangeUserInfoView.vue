<template>
    <div class="root">
    <h1>修改用户信息</h1>

    <br/>
    <el-upload
    class="avatar-uploader"
    action="http://admin.abtxw.com/api/changeAvatarFile"
    :show-file-list="false"
    :on-success="handleAvatarSuccess"
    :before-upload="beforeAvatarUpload">
      <img v-if="imageUrl" :src="imageUrl" class="avatar">
      <i v-else class="el-icon-plus avatar-uploader-icon"></i>
    </el-upload>

    <el-form class="register-container" ref="changeForm" :model="changeForm" label-width="80px">
      <el-form-item label="用户名">
        <el-input v-model="changeForm.name" :value="changeForm.name"></el-input>
      </el-form-item>

      <el-form-item label="电子邮箱">
        <el-input v-model="changeForm.email" :value="changeForm.email"></el-input>
      </el-form-item>

      <el-form-item label="密码">
        <el-input v-model="changeForm.password" :value="changeForm.password"></el-input>
      </el-form-item>

      <el-form-item label="手机号">
        <el-input type="number" v-model="changeForm.phoneNumber" :value="changeForm.phoneNumber"></el-input>
      </el-form-item>

      <el-form-item label="备注">
        <el-input v-model="changeForm.remarks" :value="changeForm.remarks"></el-input>
      </el-form-item>

      <div>
        <el-button  type="primary" style="width:100%;background-color: #505458;border: none;" @click="ConfirmChange">确认修改</el-button>
      </div>

      <br/>

      <div>
      <el-button  type="primary" style="width:100%;background-color: #505458;border: none;" @click="Cancel">取消</el-button>
      </div>
      
    </el-form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
    name: 'ChangeUserInfoView',
    data(){
        return {
            imageUrl: '',
            avatarURL: '',
            changeInfo: {
                adminName: '', //存储管理员名
                userName: '', //存储用户名
            },
            changeForm: {
                name: '',
                email: '',
                password: '',
                phoneNumber: '',
                remarks: '',
                adminName: '',
                userAvatarName: '',
                originalUserName: '',
                loginPrivilege: '',
            },
            originalUserData: {
                name: '',
                email: '',
                password: '',
                phoneNumber: '',
                remarks: '',
            }
        }
    },

    created(){
    // 获取路由参数中的用户名
    this.changeForm.adminName = this.$route.query.loginName || '';
    this.changeInfo.adminName = this.$route.query.loginName || '';
    this.changeInfo.userName = this.$route.query.userName || '';

    //获取该用户的原始数据并显示在网页中
    this.getOriginalUserData();

    //显示初始头像
    this.locateUserAvatar();

    //获取初始头像文件名
    this.getOriginalAvatarName();

    //获取用户登录权限
    this.getLoginPrivilege();
    },

    methods: {
        async locateUserAvatar() {
            try {
          const response = await axios.post('http://admin.abtxw.com/api/locateUserAvatar',this.changeInfo,{
            headers: {
              'Content-Type': 'application/json',
            },
            responseType: 'blob', // 告诉axios返回Blob对象
          });

            //处理后端响应
            const fileData = response.data; // 假设后端返回的是文件数据
            const objectURL = URL.createObjectURL(fileData); // 创建对象 URL
            this.avatarURL = objectURL

            // 一旦不再需要对象 URL，记得释放它以防止内存泄漏
            // 可以在组件销毁时或不再需要对象 URL 时调用
            // URL.revokeObjectURL(objectURL);
            } catch (error) {
            //处理错误
            console.error(error);
            //将错误抛出
            throw error;
            }
            this.imageUrl = this.avatarURL;
        },

        async getOriginalAvatarName() {
            try {
                const response = await axios.post('http://admin.abtxw.com/api/getOriginalAvatarName',this.changeInfo,{
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });

                //处理后端响应
                this.changeForm.userAvatarName = response.data.originalAvatarName;
            } catch (error) {
                console.error(error);
                throw error;
            }
        },

        async getLoginPrivilege() {
          try {
                const response = await axios.post('http://admin.abtxw.com/api/getLoginPrivilege',this.changeInfo,{
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });

                //处理后端响应
                this.changeForm.loginPrivilege = response.data.loginPrivilege;
            } catch (error) {
                console.error(error);
                throw error;
            }
        },

        async getOriginalUserData() {
            try {
                const response = await axios.post('http://admin.abtxw.com/api/getOriginalUserData',this.changeInfo,{
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });

                //处理后端响应
                this.originalUserData = response.data; // 后端返回的数据是一个数组，将其赋值给 originalUserData
            } catch (error) {
                console.error(error);
                throw error;
            }

            
            this.changeForm.name = this.originalUserData.name
            this.changeForm.email = this.originalUserData.email
            this.changeForm.password = this.originalUserData.password
            this.changeForm.phoneNumber = this.originalUserData.phoneNumber
            this.changeForm.originalUserName = this.originalUserData.name
            this.changeForm.remarks = this.originalUserData.remarks
        },

        async handleAvatarSuccess(res, file) {
            this.imageUrl = URL.createObjectURL(file.raw);

            axios.get('http://admin.abtxw.com/api/getAvatarFilename')
                .then((response) => {
                    const filename = response.data.filename;
                    this.changeForm.userAvatarName = filename;
                })
                .catch((error) => {
                    console.error('Error:', error);
                });
        },

        beforeAvatarUpload(file) {
            const isJPG = file.type === 'image/jpeg';
            const isLt2M = file.size / 1024 / 1024 < 2;

            if (!isJPG) {
            this.$message.error('上传头像图片只能是 JPG 格式!');
            }
            if (!isLt2M) {
            this.$message.error('上传头像图片大小不能超过 2MB!');
            }
            return isJPG && isLt2M;
        },

        async ConfirmChange() {
            try {
                if (
                this.changeForm.name &&
                this.changeForm.email &&
                this.changeForm.password &&
                this.changeForm.phoneNumber
                ){
                //调用函数将数据发送到后端
                await this.sendDataToBackend();

                this.$message({
                    message: '修改成功，欢迎回到用户管理系统',
                    type: 'success'
                });

                this.$router.push({
                    path:'/Home',
                    query: {
                    loginName: this.changeInfo.adminName,
                    }
                });
                } else {
                    this.$message.error('请填写完整的修改信息');
                }
            } catch (error) {
                this.$message.error('修改失败，请重新注册或联系管理员');
                console.error(error);
            }
        },

        async sendDataToBackend() {
        try {
          const response = await axios.post('http://admin.abtxw.com/api/changeData',this.changeForm,{
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
        },

        Cancel() {
          this.$router.push({
                    path:'/Home',
                    query: {
                    loginName: this.changeInfo.adminName,
                    }
                });
        },
    }
}
</script>

<style>
.root {
  text-align: center;
}
.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.avatar-uploader .el-upload:hover {
  border-color: #409EFF;
}
.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
}
.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
.register-container {
        margin: 20px auto;
        width:350px;
    }
</style>