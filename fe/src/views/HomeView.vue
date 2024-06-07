<template>
  <div class="home">
    <br/>
    <h1>用户管理系统</h1>
    
    <br/><br/>

    <div class="block" @click="toChangeAdminInfo" style="text-align:center">
      <el-avatar :size="100" :src="avatarURL" style="cursor: pointer;"></el-avatar>
    </div>

    <div class="search-container">
      <el-input placeholder="请输入用户名进行检索" v-model="searchInput" class="input-with-select">
        <el-button slot="append" icon="el-icon-search" @click="searchUser"></el-button>
      </el-input>
    </div>

    <div class="homeContainer">
      <el-table
    :data="tableData"
    style="width: 100%"
    max-height="500"
    >

      <el-table-column
        prop="name"
        label="用户名"
        width="120">
        <template slot-scope="scope">
          <span :class="{ 'highlighted': scope.row.highlighted }">{{ scope.row.name }}</span>
        </template>
      </el-table-column>

      <el-table-column
        prop="email"
        label="电子邮箱"
        width="200">
      </el-table-column>

      <el-table-column
        prop="password"
        label="密码"
        width="150">
      </el-table-column>

      <el-table-column
        prop="phoneNumber"
        label="手机号"
        width="150">
      </el-table-column>

      <el-table-column
        prop="remarks"
        label="备注"
        width="90">
      </el-table-column>

      <el-table-column
        prop="loginPrivilege"
        label="允许登录"
        width="90">
      </el-table-column>

      <el-table-column fixed="right" prop="change" label="查看用户信息" width="120">
        <template slot-scope="scope">
          <el-button type="primary" @click="toCheckUserInfo(scope.$index, tableData)">
            查看信息
          </el-button>
        </template>
      </el-table-column>

      <el-table-column fixed="right" prop="change" label="更改用户信息" width="120">
        <template slot-scope="scope">
          <el-button type="warning" @click="toChangeUserInfo(scope.$index, tableData)">
            更改信息
          </el-button>
        </template>
      </el-table-column>

      <el-table-column fixed="right" label="改变登陆权限" width="120">
        <template slot-scope="scope">
          <el-button @click="changeLoginPrivilege(scope.$index, tableData)" type="danger">
            更改权限
          </el-button>
        </template>
      </el-table-column>

      <el-table-column fixed="right" label="删除用户" width="120">
        <template slot-scope="scope">
          <el-button @click="deleteRow(scope.$index, tableData)" type="danger">
            移除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    </div>
    <el-button round @click="addUser">添加用户</el-button>
    <el-button round @click="logout">退出登录</el-button>
    <el-button round @click="deleteAdmin" type="danger" title="管理员信息及用户列表将被删除！">注销管理员账户</el-button>
    <el-button type="info" round @click="testUserLogin">测试用户登录</el-button>
  </div>
</template>

<script>
import axios from 'axios'
import { Image } from 'element-ui';

export default {
    name: 'HomeView',
    
    data() {
        return {
            circleUrl: "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png",
            squareUrl: "https://cube.elemecdn.com/9/c2/f0ee8a3c7c9638a54940382568c9dpng.png",
            avatarURL: '',
            userAvatarURL: [],
            searchInput: '',
            adminInfo: {
                loginName: '', //存储用户名
            },
            deleteName: {
                userName: '',
                adminName: '',
            },
            tableData: [],
        };
    },
    created() {
        // 获取路由参数中的用户名
        this.adminInfo.loginName = this.$route.query.loginName || '';
        this.fetchUserData();
        this.locateAvatar();
    },
    methods: {
        async fetchUserData() {
            try {
                const response = await axios.post('http://localhost:6230/api/fetchUserData', this.adminInfo, {
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                //处理后端响应
                this.tableData = response.data; // 后端返回的数据是一个数组，将其赋值给 tableData
                // this.getUserAvatarURL();
            }
            catch (error) {
                console.error(error);
                throw error;
            }
        },

        // async getUserAvatarURL() {
        //     for (let i = 0; i < this.tableData.length; i++) {
        //         try {
        //             const response = await axios.post('http://localhost:6230/api/locateHomeUserAvatar', this.tableData[i], {
        //                 headers: {
        //                     'Content-Type': 'application/json',
        //                 },
        //                 responseType: 'blob', // 告诉axios返回Blob对象
        //             });
        //             //处理后端响应
        //             const fileData = response.data; // 假设后端返回的是文件数据
        //             const objectURL = URL.createObjectURL(fileData); // 创建对象 URL
        //             this.userAvatarURL.push(objectURL);
        //             // 一旦不再需要对象 URL，记得释放它以防止内存泄漏
        //             // 可以在组件销毁时或不再需要对象 URL 时调用
        //             // URL.revokeObjectURL(objectURL);
        //         }
        //         catch (error) {
        //             //处理错误
        //             console.error(error);
        //             //将错误抛出
        //             throw error;
        //         }
        //     }
        // },

        async locateAvatar() {
            try {
                const response = await axios.post('http://localhost:6230/api/locateAvatar', this.adminInfo, {
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    responseType: 'blob', // 告诉axios返回Blob对象
                });
                //处理后端响应
                const fileData = response.data; // 假设后端返回的是文件数据
                const objectURL = URL.createObjectURL(fileData); // 创建对象 URL
                this.avatarURL = objectURL;
                // 一旦不再需要对象 URL，记得释放它以防止内存泄漏
                // 可以在组件销毁时或不再需要对象 URL 时调用
                // URL.revokeObjectURL(objectURL);
            }
            catch (error) {
                //处理错误
                console.error(error);
                //将错误抛出
                throw error;
            }
        },

        toChangeAdminInfo() {
          this.$router.push({
                path: '/ChangeAdminInfoView',
                query: {
                    loginName: this.adminInfo.loginName
                }
            });
        },

        async deleteRow(index, rows) {
            this.deleteName.adminName = this.adminInfo.loginName;
            this.deleteName.userName = rows[index].name;
            try {
                const response = await axios.post('http://localhost:6230/api/deleteUserByName', this.deleteName, {
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                //处理后端响应
                console.log("message");
            }
            catch (error) {
                //处理错误
                console.error(error);
                //将错误抛出
                throw error;
            }
            this.fetchUserData();
        },

        async changeLoginPrivilege(index, rows) {
            this.deleteName.adminName = this.adminInfo.loginName;
            this.deleteName.userName = rows[index].name;
            try {
                const response = await axios.post('http://localhost:6230/api/changeLoginPrivilege', this.deleteName, {
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                //处理后端响应
                console.log("message");
            }
            catch (error) {
                //处理错误
                console.error(error);
                //将错误抛出
                throw error;
            }
            this.fetchUserData();
        },

        toChangeUserInfo(index, rows) {
            this.$router.push({
                path: '/ChangeUserInfoView',
                query: {
                    loginName: this.adminInfo.loginName,
                    userName: rows[index].name
                }
            });
        },

        addUser() {
            this.$router.push({
                path: '/AddUserView',
                query: {
                    loginName: this.adminInfo.loginName
                }
            });
        },

        logout() {
            this.$router.push({ path: '/' });
        },

        searchUser() {
            // 获取输入框的值
            const searchTerm = this.searchInput.trim();

            if (searchTerm == 0) {
              this.$message.error('请在输入框内输入用户名！');
            } else {
              // 遍历表格数据，查找匹配的用户
              let i = 0; //用于判断是否匹配成功
              this.tableData.forEach((user, index) => {
                  if (user.name === searchTerm) {
                    this.toShowUserInfo(this.tableData[index].name)

                    i = 1;
                  }
              });

              if ( i == 0) {
                this.$message.error('该用户不存在！');
              }
            }
        },

        async deleteAdmin() {
            try {
                const response = await axios.post('http://localhost:6230/api/deleteAdmin', this.adminInfo, {
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                //处理后端响应
                console.log("message"); // 后端返回的数据是一个数组，将其赋值给 tableData
                //路由回登陆页面
                this.$router.push({ path: '/' });
            }
            catch (error) {
                console.error(error);
                throw error;
            }
        },

        async toCheckUserInfo(index, rows) {
          this.$router.push({
                path: '/CheckUserInfoView',
                query: {
                    loginName: this.adminInfo.loginName,
                    userName: rows[index].name
                }
            });
        },

        async toShowUserInfo(name) {
          this.$router.push({
                path: '/ChangeUserInfoView',
                query: {
                    loginName: this.adminInfo.loginName,
                    userName: name
                }
            });
        },

        async testUserLogin() {
          this.$router.push({
                path: '/TestUserLoginView',
                query: {
                    loginName: this.adminInfo.loginName,
                }
            });
        },
    },
    
}
</script>

<style>
  .home {
    text-align: center;
  }
  .homeContainer {
    border-radius: 15px;
    background-clip: padding-box;
    margin: 60px auto;
    width:1400px;
    padding: 28px;
    background: #fff;
    border:1px solid #eaeaee;
    box-shadow: 0 0 25px #cac6c6;
  }
  .search-container {
    display: inline-block;
    width: 500px;
    margin-top: 15px;
  }
  .highlighted {
    background-color: #99e1f3; /* 更改为您想要的颜色 */
  }
</style>