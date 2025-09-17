// frontend/src/App.js
import React, { useEffect, useState } from "react";
import Recommend from "./Recommend";
import './App.css';
import {
  fetchCourses,
  fetchStudents,
  fetchStudentCourses,
  createStudent,
} from "./api";

export default function App() {
  // 视图：home / recommend
  const [view, setView] = useState("home");

  // ===== 下面是你原来的 Home 页面状态与逻辑 =====
  const [students, setStudents] = useState([]);
  const [courses, setCourses] = useState([]);
  const [count, setCount] = useState(0);
  const [selectedStudent, setSelectedStudent] = useState("");
  const [newStudentName, setNewStudentName] = useState("");

  const handleAddStudent = async () => {
    if (!newStudentName.trim()) return;
    try {
      await createStudent({ name: newStudentName.trim() });
      setNewStudentName("");
      fetchStudents().then((data) => setStudents(data.items || [])); // 刷新学生列表
    } catch (err) {
      alert("新增学生失败: " + err.message);
    }
  };

  const loadAllCourses = async () => {
    const data = await fetchCourses();
    setCourses(data.items || []);
    setCount(data.count || 0);
  };

  useEffect(() => {
    // 初始加载全部课程与学生列表
    loadAllCourses();
    fetchStudents().then((data) => setStudents(data.items || []));
  }, []);

  const onSelectStudent = async (e) => {
    const id = e.target.value;
    setSelectedStudent(id);
    if (!id) {
      await loadAllCourses();
      return;
    }
    const data = await fetchStudentCourses(id);
    setCourses(data.items || []);
    setCount(data.count || 0);
  };

  // ====== 视图渲染 ======
  return (
    <div className="container" style={{ padding: 16 }}>
      {/* 顶部导航（在同一页面切换视图） */}
      <div style={{ marginBottom: 16 }}>
        <button
          onClick={() => setView("home")}
          style={{ marginRight: 8 }}
          disabled={view === "home"}
        >
          首页
        </button>
        <button
          onClick={() => setView("recommend")}
          disabled={view === "recommend"}
        >
          AI 推荐选课
        </button>
      </div>

      {view === "recommend" ? (
        <Recommend />
      ) : (
        // ======= 你的原始“首页”界面 =======
        <div>
          <h1>RIC 选课平台</h1>

          <div className="toolbar" style={{ marginBottom: 12 }}>
            <button onClick={loadAllCourses}>查看所有课程</button>

            <label style={{ marginLeft: 12 }}>
              选择学生：
              <select value={selectedStudent} onChange={onSelectStudent}>
                <option value="">（全部）</option>
                {students.map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.name}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div className="summary" style={{ marginBottom: 8 }}>
            共 <strong>{count}</strong> 门课程
          </div>

          <table>
            <thead>
              <tr>
                <th>课程代码</th>
                <th>课程名称</th>
              </tr>
            </thead>
            <tbody>
              {courses.map((c) => (
                <tr key={c.id}>
                  <td>{c.code}</td>
                  <td>{c.title}</td>
                </tr>
              ))}
            </tbody>
          </table>

          <div className="add-student" style={{ marginTop: 16 }}>
            <h2>新增学生</h2>
            <input
              type="text"
              placeholder="学生姓名"
              value={newStudentName}
              onChange={(e) => setNewStudentName(e.target.value)}
            />
            <button onClick={handleAddStudent} style={{ marginLeft: 8 }}>
              添加
            </button>
          </div>

          <footer style={{ marginTop: 16 }}>
            <small>前端：React | 后端：Go | 数据库：PostgreSQL</small>
          </footer>
        </div>
      )}
    </div>
  );
}
