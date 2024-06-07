<template>
    <div class="root">
    <h1>添加用户</h1>

    <br/>
    <el-upload
    class="avatar-uploader"
    action="http://localhost:6230/api/uploadAvatarFile"
    :show-file-list="false"
    :on-success="handleAvatarSuccess"
    :before-upload="beforeAvatarUpload">
      <img v-if="imageUrl" :src="imageUrl" class="avatar">
      <i v-else class="el-icon-plus avatar-uploader-icon"></i>
    </el-upload>

    <el-form class="register-container" ref="addForm" :model="addForm" label-width="80px">
      <el-form-item label="用户名">
        <el-input v-model="addForm.name"></el-input>
      </el-form-item>

      <el-form-item label="电子邮箱">
        <el-input v-model="addForm.email"></el-input>
      </el-form-item>

      <el-form-item label="密码">
        <el-input v-model="addForm.password"></el-input>
      </el-form-item>

      <el-form-item label="确认密码">
          <el-input v-model="confirmPassword" @input="confirmPasswordFunc"></el-input>
          <h5 style="color: red;">{{ confirmWarn }}</h5>
      </el-form-item>

      <el-form-item label="手机号">
        <el-input type="number" v-model="addForm.phoneNumber"></el-input>
      </el-form-item>

      <el-form-item label="备注">
        <el-input v-model="addForm.remarks"></el-input>
      </el-form-item>

      <div>
        <el-button  type="primary" style="width:100%;background-color: #505458;border: none;" @click="AddUser">确认添加</el-button>
      </div>

      <br/>

      <div>
        <el-button  type="primary" style="width:100%;background-color: #505458;border: none;" @click="Cancel">取消</el-button>
      </div>
      
    </el-form>
  </div>
</template>

<script>
import axios from 'axios'

export default {
    name: 'AddUserView',

    data() {
        return {
        imageUrl: '',
        confirmPassword: '',
        confirmWarn: '',
        addForm: {
            name: '',
            email: '',
            password: '',
            phoneNumber: '',
            avatarFilename: '',
            remarks: '',
            adminName: '',
            },
        }
    },

    created() {
        // 获取路由参数中的用户名
        this.addForm.adminName = this.$route.query.loginName || '';
    },

    methods: {
        handleAvatarSuccess(res, file) {
        this.imageUrl = URL.createObjectURL(file.raw);

        axios.get('http://localhost:6230/api/getAvatarFilename')
                .then((response) => {
                    const filename = response.data.filename;
                    this.addForm.avatarFilename = filename;
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

        confirmPasswordFunc() {
        if(this.addForm.password === this.confirmPassword) {
          this.confirmWarn = '';
          return;
        } else {
          this.confirmWarn = '密码不一致，请从新填写！';
          return;
        }
      },

      async AddUser() {
        try {
          if (
          this.addForm.name &&
          this.addForm.email &&
          this.addForm.password &&
          this.addForm.phoneNumber
          ){
          //调用函数将数据发送到后端
          await this.sendDataToBackend();

          this.$message({
            message: '添加成功，欢迎回到用户管理系统',
            type: 'success'
          });

          this.$router.push({
            path:'/Home',
            query: {
              loginName: this.addForm.adminName,
            }
          });
          } else {
            this.$message.error('请填写完整的要添加的用户信息');
          }
        } catch (error) {
          this.$message.error('添加失败，请重新添加或联系管理员');
          console.error(error);
        }
      },

      async sendDataToBackend() {
        try {
          const response = await axios.post('http://localhost:6230/api/addUser',this.addForm,{
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
              loginName: this.addForm.adminName,
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