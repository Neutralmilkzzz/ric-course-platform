// frontend/src/Recommend.js
import React, { useEffect, useState } from "react";
import { fetchCourses } from "./api";
import { getRecommendations } from "./deepseekApi";
import './App.css';
export default function Recommend() {
  const [allCourseTitles, setAllCourseTitles] = useState([]);
  const [major, setMajor] = useState("");
  const [loading, setLoading] = useState(false);
  const [answer, setAnswer] = useState("");

  useEffect(() => {
    // 加载全部课程标题
    (async () => {
      try {
        const data = await fetchCourses();
        const titles = (data.items || []).map((c) => c.title);
        setAllCourseTitles(titles);
      } catch (e) {
        alert("加载课程失败：" + e.message);
      }
    })();
  }, []);

  const handleRecommend = async () => {
    if (!major.trim()) {
      alert("请先输入专业");
      return;
    }
    if (!allCourseTitles.length) {
      alert("课程列表为空，无法生成推荐");
      return;
    }
    setLoading(true);
    setAnswer("");
    try {
      const text = await getRecommendations(allCourseTitles, major.trim());
      setAnswer(text);
    } catch (e) {
      alert("生成推荐失败：" + e.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="recommend-container" style={{ padding: 16 }}>
      <h2>AI 推荐选课</h2>

      <div style={{ marginBottom: 12 }}>
        <label>
          专业：
          <input
            type="text"
            placeholder="例如：Computer Science / Finance / EE"
            value={major}
            onChange={(e) => setMajor(e.target.value)}
            style={{ marginLeft: 8 }}
          />
        </label>
        <button
          onClick={handleRecommend}
          disabled={loading}
          style={{ marginLeft: 12 }}
        >
          {loading ? "生成中..." : "生成推荐"}
        </button>
      </div>

      <div style={{ marginBottom: 8, fontSize: 13, color: "#666" }}>
        已加载课程数量：{allCourseTitles.length}
      </div>

      {answer && (
        <div
          style={{
            whiteSpace: "pre-wrap",
            background: "#243057",
            padding: 12,
            borderRadius: 8,
          }}
        >
          {answer}
        </div>
      )}
    </div>
  );
}
