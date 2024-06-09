<template>
    <div class="root">
    <h1>查看用户信息</h1>

    <br/>
    <el-upload
    class="avatar-uploader"
    action="http://admin.abtxw.com/api/changeAvatarFile"
    :show-file-list="false"
    :on-success="handleAvatarSuccess"
    :before-upload="beforeAvatarUpload"
    :disabled="true">
      <img v-if="imageUrl" :src="imageUrl" class="avatar">
      <i v-else class="el-icon-plus avatar-uploader-icon"></i>
    </el-upload>

    <el-form class="register-container" ref="checkForm" :model="checkForm" label-width="80px">
      <el-form-item label="用户名">
        <el-input v-model="checkForm.name" :value="checkForm.name"></el-input>
      </el-form-item>

      <el-form-item label="电子邮箱">
        <el-input v-model="checkForm.email" :value="checkForm.email"></el-input>
      </el-form-item>

      <el-form-item label="密码">
        <el-input v-model="checkForm.password" :value="checkForm.password"></el-input>
      </el-form-item>

      <el-form-item label="手机号">
        <el-input type="number" v-model="checkForm.phoneNumber" :value="checkForm.phoneNumber"></el-input>
      </el-form-item>

      <el-form-item label="备注">
        <el-input v-model="checkForm.remarks" :value="checkForm.remarks"></el-input>
      </el-form-item>

      <div>
      <el-button  type="primary" style="width:100%;background-color: #505458;border: none;" @click="goBack">返回</el-button>
      </div>
      
    </el-form>
  </div>
</template>

<script>
import axios from 'axios'

export default {
    name: 'CheckUserInfoView',

    data() {
        return {
            imageUrl: '',
            avatarURL: '',
            checkInfo: {
                adminName: '', //存储管理员名
                userName: '', //存储用户名
            },
            checkForm: {
                name: '',
                email: '',
                password: '',
                phoneNumber: '',
                remarks: '',
                adminName: '',
                userAvatarName: '',
                originalUserName: '',
            },
            originalAdminData: {
                name: '',
                email: '',
                password: '',
                phoneNumber: '',
                remarks: '',
            }
        }
    },

    created() {
            this.checkForm.adminName = this.$route.query.loginName || '';
            this.checkInfo.adminName = this.checkForm.adminName || '';
            this.checkInfo.userName = this.$route.query.userName || '';

            //获取该用户的原始数据并显示在网页中
            this.getOriginalUserData();

            //显示初始头像
            this.locateUserAvatar();
        },

        methods: {
        async locateUserAvatar() {
            try {
          const response = await axios.post('http://admin.abtxw.com/api/locateUserAvatar',this.checkInfo,{
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

        async getOriginalUserData() {
            try {
                const response = await axios.post('http://admin.abtxw.com/api/getOriginalUserData',this.checkInfo,{
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

            
            this.checkForm.name = this.originalUserData.name
            this.checkForm.email = this.originalUserData.email
            this.checkForm.password = this.originalUserData.password
            this.checkForm.phoneNumber = this.originalUserData.phoneNumber
            this.checkForm.originalUserName = this.originalUserData.name
            this.checkForm.remarks = this.originalUserData.remarks
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

        goBack() {
          this.$router.push({
                    path:'/Home',
                    query: {
                    loginName: this.checkInfo.adminName,
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
  cursor: default;
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