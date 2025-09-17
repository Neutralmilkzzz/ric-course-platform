// frontend/src/deepseekApi.js
const DEEPSEEK_BASE =
  process.env.REACT_APP_DEEPSEEK_BASE_URL || "https://api.deepseek.com";
const DEEPSEEK_KEY = process.env.REACT_APP_DEEPSEEK_KEY; // 在 .env 里配置

/**
 * 调用 DeepSeek（OpenAI 兼容格式）获取推荐
 * @param {string[]} courses - 课程标题数组
 * @param {string} major - 专业名称
 * @returns {Promise<string>} - 模型返回的文本
 */
export async function getRecommendations(courses, major) {
  if (!DEEPSEEK_KEY) {
    throw new Error(
      "缺少 REACT_APP_DEEPSEEK_KEY，请在 frontend/.env 中配置后重建前端"
    );
  }

  const systemPrompt =
    `你是大学选课顾问。以下是可选课程列表（仅供参考）：\n` +
    courses.map((t, i) => `${i + 1}. ${t}`).join("\n") +
    `\n请根据学生专业「${major}」给出推荐课程，并简要说明理由。输出尽量精炼，给出 3-5 门即可。`;

  const url = `${DEEPSEEK_BASE}/v1/chat/completions`;
  const payload = {
    model: "deepseek-chat", // 如有不同型号，改这里
    messages: [
      { role: "system", content: systemPrompt },
      {
        role: "user",
        content:
          "请基于上面课程与专业，推荐我本学期该选择的课程，并说明理由。",
      },
    ],
    temperature: 0.7,
  };

  const res = await fetch(url, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${DEEPSEEK_KEY}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`DeepSeek API 调用失败：${res.status} ${text}`);
  }

  const data = await res.json();
  const content =
    data?.choices?.[0]?.message?.content ??
    "(没有返回内容，检查模型/配额/参数)";
  return content;
}
