const BASE = process.env.REACT_APP_API_BASE_URL || "https://ric-course-platform-3.onrender.com/";

export async function fetchCourses() {
  const res = await fetch(`${BASE}/api/courses`);
  if (!res.ok) throw new Error('Failed to fetch courses');
  return res.json();
}

export async function fetchStudents() {
  const res = await fetch(`${BASE}/api/students`);
  if (!res.ok) throw new Error('Failed to fetch students');
  return res.json();
}

export async function fetchStudentCourses(id) {
  const res = await fetch(`${BASE}/api/students/${id}/courses`);
  if (!res.ok) throw new Error('Failed to fetch student courses');
  return res.json();
}

// 添加新课程
export async function createCourse(course) {
  const res = await fetch(`${BASE}/api/courses`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(course),
  });
  if (!res.ok) throw new Error("Failed to create course");
  return res.json();
}

// 更新课程
export async function updateCourse(id, course) {
  const res = await fetch(`${BASE}/api/courses/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(course),
  });
  if (!res.ok) throw new Error("Failed to update course");
  return res.json();
}

// 删除课程
export async function deleteCourse(id) {
  const res = await fetch(`${BASE}/api/courses/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) throw new Error("Failed to delete course");
  return res.json();
}

// 添加学生
export async function createStudent(student) {
  const res = await fetch(`${BASE}/api/students`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(student),
  });
  if (!res.ok) throw new Error("Failed to create student");
  return res.json();
}

// 更新学生
export async function updateStudent(id, student) {
  const res = await fetch(`${BASE}/api/students/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(student),
  });
  if (!res.ok) throw new Error("Failed to update student");
  return res.json();
}

// 删除学生
export async function deleteStudent(id) {
  const res = await fetch(`${BASE}/api/students/${id}`, {
    method: "DELETE",
  });
  if (!res.ok) throw new Error("Failed to delete student");
  return res.json();
}

