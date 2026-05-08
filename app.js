const STORAGE_KEY = "stock-portfolio-desk-v2";

const seedState = {
  totalCapital: 1150000,
  cash: 477238.13560642325,
  fx: { CNY: 1, HKD: 0.8716, USD: 7.1 },
  trades: [],
  decisionLogs: [],
  holdings: [
    {
      symbol: "0700.HK",
      name: "腾讯控股",
      shares: 200,
      cost: 480.43,
      currentPrice: 463,
      previousClose: 463,
      action: "继续持有；新资金暂不追买，放入核心替补",
      status: "未达标（安全边际<15%）",
      marginOfSafety: 0.09,
      qualityScore: 89,
      risk: "无立即否决；政策/地缘/AI投入需折价",
      industry: "互联网平台/游戏/广告/金融科技",
      currency: "HKD",
      intrinsicValue: 508,
      fairValueRange: "HK$480-560",
      targetBuyPrice: 432,
      businessModel: 28,
      moat: 23,
      governance: 17,
      financialQuality: 21,
      updatedAt: "2026-05-06；最新价约HK$463.00；HKD/CNY约0.8716；FY2025",
      notes: "FY2025：收入RMB7518亿、Non-IFRS净利RMB2596亿、FCF RMB1826亿、净现金RMB1071亿。"
    },
    {
      symbol: "000333.SZ",
      name: "美的集团",
      shares: 600,
      cost: 79.638,
      currentPrice: 80.44,
      previousClose: 80.44,
      action: "放入核心替补；A股暂不追买，H股优先但等待≤HK$86-87",
      status: "未达标（A股安全边际<20%；H股接近达标）",
      marginOfSafety: 0.153,
      qualityScore: 88,
      risk: "无立即否决；Q1扣非下滑、海外关税/汇率、价格战需跟踪",
      industry: "家电/全球化制造/ToB楼宇科技/机器人自动化",
      currency: "CNY",
      intrinsicValue: 95,
      fairValueRange: "¥90-100",
      targetBuyPrice: 76,
      businessModel: 28,
      moat: 23,
      governance: 18,
      financialQuality: 19,
      updatedAt: "2026-05-06；A股最新价约¥80.44；H股约HK$87.70；FY2025/2026Q1",
      notes: "FY2025：营收RMB4585亿、归母净利RMB439.45亿、年度分红¥4.30/股。"
    },
    {
      symbol: "002415.SZ",
      name: "海康威视",
      shares: 1200,
      cost: 34.54,
      currentPrice: 36.29,
      previousClose: 36.29,
      action: "重点预期差候选/核心替补边缘；可小仓验证，不宜重仓；Q2验证后再升级",
      status: "未达标（安全边际约13.6%<25%；预期差仓可观察）",
      marginOfSafety: 0.136,
      qualityScore: 84,
      risk: "无一票否决；地缘/合规/实体清单、Q1经营现金流为负、AIoT重估需验证",
      industry: "AIoT/安防/机器视觉/科技制造平台",
      currency: "CNY",
      intrinsicValue: 42,
      fairValueRange: "¥34-48",
      targetBuyPrice: 31.5,
      businessModel: 25,
      moat: 23,
      governance: 16,
      financialQuality: 20,
      updatedAt: "2026-05-06；最新价约¥36.29；FY2025/2026Q1；董秘大额增持后修正",
      notes: "FY2025：营收约RMB925.08亿、归母净利约RMB141.95亿；2026Q1归母净利同比+36.42%。"
    },
    {
      symbol: "600887.SH",
      name: "伊利股份",
      shares: 1300,
      cost: 26.469,
      currentPrice: 27.45,
      previousClose: 27.45,
      action: "放入核心替补；暂不追买，等待¥24-26",
      status: "未达标（安全边际约14.2%<25%）",
      marginOfSafety: 0.1421875,
      qualityScore: 83,
      risk: "无一票否决；需求弱复苏、原奶上涨传导不顺、液奶仍下滑、食品安全风险需跟踪",
      industry: "乳制品/消费龙头/高股息/奶周期修复",
      currency: "CNY",
      intrinsicValue: 32,
      fairValueRange: "¥28-36",
      targetBuyPrice: 24,
      businessModel: 24,
      moat: 22,
      governance: 16,
      financialQuality: 21,
      updatedAt: "2026-05-07；最新价约¥27.45；FY2025/2026Q1；奶周期底部右侧观察",
      notes: "2025拟派息¥1.38/股，按¥27.45股息率约5.0%；达标买入价≤¥24。"
    },
    { symbol: "600036.SH", name: "招商银行", shares: 500, cost: 39.18, currentPrice: 39.18, previousClose: 39.18, currency: "CNY", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "0696.HK", name: "民航信", shares: 11000, cost: 10.648, currentPrice: 10.648, previousClose: 10.648, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "0506.HK", name: "中国食品", shares: 22000, cost: 4.041, currentPrice: 4.041, previousClose: 4.041, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "2669.HK", name: "中海物业", shares: 20000, cost: 4.468, currentPrice: 4.468, previousClose: 4.468, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "6049.HK", name: "保利物业", shares: 2600, cost: 32.663, currentPrice: 32.663, previousClose: 32.663, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "0883.HK", name: "中海油", shares: 2000, cost: 29.326, currentPrice: 29.326, previousClose: 29.326, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    {
      symbol: "1448.HK",
      name: "福寿园",
      shares: 11000,
      cost: 2.521,
      currentPrice: 2.64,
      previousClose: 2.64,
      currentPriceDate: "2026-05-07",
      previousCloseDate: "2026-05-07",
      action: "暂不行动；不买入；不纳入核心替补，等待2025年报、审计意见、法证调查结论和复牌后再重估",
      status: "未达标（停牌、年报延迟、治理与财务可靠性风险未解除）",
      marginOfSafety: 0,
      qualityScore: 62,
      risk: "已触发重大风险否决项：停牌、业绩延迟、现金及采购付款事项调查、管理层/内控可信度下降、墓穴ASP大幅下滑、资产和商誉减值风险",
      industry: "殡葬服务/墓园运营/生命服务",
      currency: "HKD",
      intrinsicValue: 2.65,
      fairValueRange: "HK$1.6-3.1",
      targetBuyPrice: 2,
      businessModel: 22,
      moat: 16,
      governance: 5,
      financialQuality: 19,
      updatedAt: "2026-05-07；停牌前最后价约HK$2.64；用户更新分析",
      notes: "计划：剔除/仅风险观察。复牌前不行动；复牌后若审计无保留、调查无重大重述且价格≤HK$2.0-2.2，才重新评估普通候选价值。纪律：质量分低于75且有重大风险否决项；不因低估值或净现金买入，先等风险解除。最新市场状态：股份自2026-03-20起停牌，停牌前最后价约HK$2.64。最新可用财务口径：2024收入约RMB20.77亿，归母净利约RMB3.73亿，EPS约RMB0.164；2025H1收入约RMB6.11亿，归母亏损约RMB2.61亿，EPS约-RMB0.115。核心判断：福寿园当前不是单纯估值杀，而是业绩杀、治理杀和财报可信度风险叠加；内在价值区间仅为压力测试，不作为可执行买入依据。"
    },
    {
      symbol: "07489.HK",
      name: "岚图汽车",
      shares: 2132,
      cost: 0,
      currentPrice: 5.89,
      previousClose: 5.89,
      currentPriceDate: "2026-05-07",
      previousCloseDate: "2026-05-07",
      action: "放入普通候选池观察；当前不买入，等待扣非利润和自由现金流验证",
      status: "未达标（质量分<75且安全边际不足）",
      marginOfSafety: 0.16,
      qualityScore: 72,
      risk: "盈利质量受政府补助影响，梦想家单一车型依赖较高，新能源车价格战和智能化竞争可能压缩毛利率",
      industry: "新能源乘用车/高端MPV/央企汽车",
      currency: "HKD",
      intrinsicValue: 7,
      fairValueRange: "HK$4.5-8.5",
      targetBuyPrice: 4.8,
      businessModel: 21,
      moat: 16,
      governance: 16,
      financialQuality: 19,
      updatedAt: "2026-05-07；估值基于HK$5.89附近股价；用户更新分析",
      notes: "2025年收入约人民币348.65亿元，毛利率约20.9%，净利润约人民币10.17亿元，首次年度盈利；2025年交付约150169辆，2026年1-4月交付约49038辆。估值基于HK$5.89附近股价、市值约HK$216.8亿、PE约16.9倍、PB约1.78倍。核心假设是2026年需验证扣非利润、经营现金流和自由现金流质量。"
    }
  ],
  plan: [
    { rank: 1, name: "腾讯控股", priority: "观察/低优先级", advice: "继续持有；新资金等待≤HK$432，HK$400-430可分批", discipline: "优秀资产要求≥15%安全边际；当前约9%，未达标" },
    { rank: 2, name: "美的集团", priority: "核心替补/中优先级", advice: "A股等待≤¥76分批；H股≤HK$86-87优先；当前不追买", discipline: "优秀资产要求≥20%安全边际；A股当前约15.3%，未达标" },
    { rank: 3, name: "海康威视", priority: "重点预期差候选/中优先级", advice: "不重仓；¥35-37仅适合小仓验证，¥30-32更从容；Q2验证后可升核心替补", discipline: "质量分84，合格候选要求≥25%安全边际" },
    { rank: 4, name: "伊利股份", priority: "核心替补/中低优先级", advice: "暂不追买；¥25-26开始关注，≤¥24可考虑分批", discipline: "质量分83，合格候选要求≥25%安全边际" },
    { rank: 99, name: "岚图汽车", priority: "普通候选池/低优先级", advice: "HK$4.2-4.8才接近可观察买入区；若2026H1扣非利润和自由现金流转正，可重新上修估值", discipline: "质量分低于75原则上不进入核心资产池；安全边际不足时不试仓" }
  ],
  candidates: [
    {
      symbol: "600690.SH",
      name: "海尔智家A",
      status: "候选池",
      action: "放入普通候选池观察；A股暂不追，H股赔率更优",
      marginOfSafety: 0.17,
      qualityScore: 83,
      industry: "家电/全球化白电/智慧家庭",
      currency: "CNY",
      intrinsicValue: 26,
      fairValueRange: "¥24-28",
      targetBuyPrice: 19.5
    }
  ],
  rules: [
    { dimension: "商业模式", score: 30, standard: "需求刚性、收入可重复、定价权、资本开支、行业空间" },
    { dimension: "护城河", score: 25, standard: "品牌/规模/网络效应/牌照/成本优势、份额稳定、利润率优于同行" },
    { dimension: "管理层/企业文化/治理", score: 20, standard: "长期主义、资本配置、股东回报、披露透明、少画饼" },
    { dimension: "财务质量", score: 25, standard: "ROE/ROIC、自由现金流、资产负债表、利润率、应收/存货/资本开支" }
  ]
};

const palette = ["#087f5b", "#1c4f82", "#a16207", "#7c3aed", "#be123c", "#0f766e"];
const USE_BACKEND = location.protocol === "http:" || location.protocol === "https:";
const BUY_PROXIMITY = 0.05;
const SAFETY_MARGIN_TARGET = 0.25;

let state = loadState();
let activeFilter = "all";
let searchTerm = "";
let pendingResearch = null;
let candidateSort = "priority";
const pageTitles = {
  overview: "纸牌屋",
  positions: "当前持仓",
  trades: "交易流水",
  candidates: "候选池",
  "stock-detail": "股票分析详情"
};

const elements = {
  pageTitle: document.querySelector("#pageTitle"),
  positionsBody: document.querySelector("#positionsBody"),
  tradeList: document.querySelector("#tradeList"),
  allocationChart: document.querySelector("#allocationChart"),
  overviewPlanList: document.querySelector("#overviewPlanList"),
  opportunityRadar: document.querySelector("#opportunityRadar"),
  disciplineDashboard: document.querySelector("#disciplineDashboard"),
  triggerAlerts: document.querySelector("#triggerAlerts"),
  decisionLogList: document.querySelector("#decisionLogList"),
  candidateList: document.querySelector("#candidateList"),
  candidateSort: document.querySelector("#candidateSort"),
  stockDetail: document.querySelector("#stockDetail"),
  totalValue: document.querySelector("#totalValue"),
  totalPnl: document.querySelector("#totalPnl"),
  totalPnlRate: document.querySelector("#totalPnlRate"),
  dayChange: document.querySelector("#dayChange"),
  dayChangeRate: document.querySelector("#dayChangeRate"),
  cashBalance: document.querySelector("#cashBalance"),
  positionCount: document.querySelector("#positionCount"),
  recordCount: document.querySelector("#recordCount"),
  searchInput: document.querySelector("#searchInput"),
  updateQuotesButton: document.querySelector("#updateQuotesButton"),
  quoteUpdateStatus: document.querySelector("#quoteUpdateStatus"),
  tradeDialog: document.querySelector("#tradeDialog"),
  tradeForm: document.querySelector("#tradeForm"),
  holdingDialog: document.querySelector("#holdingDialog"),
  holdingForm: document.querySelector("#holdingForm"),
  researchDialog: document.querySelector("#researchDialog"),
  researchForm: document.querySelector("#researchForm"),
  researchJSON: document.querySelector("#researchJSON"),
  researchPreview: document.querySelector("#researchPreview"),
  researchStatus: document.querySelector("#researchStatus"),
  importResearchButton: document.querySelector("#importResearch")
};

syncCash();

function loadState() {
  if (USE_BACKEND) {
    return structuredClone(seedState);
  }

  const saved = localStorage.getItem(STORAGE_KEY);
  if (!saved) return structuredClone(seedState);

  try {
    return { ...structuredClone(seedState), ...JSON.parse(saved) };
  } catch {
    return structuredClone(seedState);
  }
}

function saveState() {
  if (USE_BACKEND) return;
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
}

async function requestJSON(path, options = {}) {
  const response = await fetch(path, {
    headers: { "Content-Type": "application/json", ...(options.headers ?? {}) },
    ...options
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: "request failed" }));
    throw new Error(error.error ?? "request failed");
  }

  return response.json();
}

async function loadBackendState() {
  if (!USE_BACKEND) return false;

  try {
    state = await requestJSON("/api/state");
    localStorage.removeItem(STORAGE_KEY);
    return true;
  } catch (error) {
    console.warn("后端不可用，使用浏览器本地数据", error);
    return false;
  }
}

function fx(currencyCode) {
  return state.fx[currencyCode] ?? 1;
}

function currency(value, currencyCode = "CNY") {
  return new Intl.NumberFormat("zh-CN", {
    style: "currency",
    currency: currencyCode,
    minimumFractionDigits: 2
  }).format(value);
}

function percent(value, signed = true) {
  const prefix = signed && value >= 0 ? "+" : "";
  return `${prefix}${value.toFixed(2)}%`;
}

function finiteNumber(value) {
  const number = Number(value);
  return Number.isFinite(number) ? number : null;
}

function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max);
}

function stockHash(symbol) {
  return `#stock=${encodeURIComponent(symbol)}`;
}

function normalizeSymbol(symbol) {
  const text = String(symbol ?? "").trim().toUpperCase();
  if (text.endsWith(".HK")) {
    const code = text.slice(0, -3);
    if (/^\d+$/.test(code)) {
      return `${String(Number(code)).padStart(4, "0")}.HK`;
    }
  }
  return text;
}

function escapeHTML(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function displayText(value, fallback = "-") {
  const text = String(value ?? "").trim();
  return text || fallback;
}

function closeDateText(position) {
  if (!position.currentPriceDate && !position.previousCloseDate) return "";
  const currentDate = position.currentPriceDate || "未知";
  const previousDate = position.previousCloseDate || "未知";
  return `今收 ${currentDate} · 昨收 ${previousDate}`;
}

function calculatedMarginOfSafety(stock) {
  if (!Number.isFinite(stock?.intrinsicValue) || stock.intrinsicValue <= 0) return null;
  if (!Number.isFinite(stock?.currentPrice) || stock.currentPrice <= 0) return null;
  return (stock.intrinsicValue - stock.currentPrice) / stock.intrinsicValue;
}

function displayMarginOfSafety(stock) {
  const computed = calculatedMarginOfSafety(stock);
  const value = Number.isFinite(computed) ? computed : stock?.marginOfSafety;
  return Number.isFinite(value) ? percent(value * 100, false) : "-";
}

function cleanResearchJSON(raw) {
  return String(raw ?? "")
    .trim()
    .replace(/^```(?:json)?\s*/i, "")
    .replace(/```$/i, "")
    .trim()
    .replace(/[“”]/g, '"')
    .replace(/[‘’]/g, "'");
}

function parseResearchJSON() {
  const text = cleanResearchJSON(elements.researchJSON.value);
  if (!text) throw new Error("请粘贴 ChatGPT 生成的 JSON");
  return JSON.parse(text);
}

function setResearchStatus(message, tone = "") {
  elements.researchStatus.textContent = message;
  elements.researchStatus.className = `research-status ${tone}`.trim();
}

function setQuoteUpdateStatus(message, tone = "") {
  elements.quoteUpdateStatus.textContent = message;
  elements.quoteUpdateStatus.className = tone;
}

function appendClientDecisionLog(entry) {
  const nextEntry = {
    id: Date.now(),
    date: new Date().toISOString().slice(0, 19).replace("T", " "),
    type: "event",
    symbol: "",
    name: "",
    price: null,
    currency: "",
    decision: "",
    discipline: "",
    detail: "",
    ...entry
  };
  state.decisionLogs = [...(state.decisionLogs ?? []), nextEntry].slice(-500);
}

function targetTypeText(type) {
  if (type === "holding") return "更新现有持仓";
  if (type === "candidate") return "更新候选池";
  if (type === "newCandidate") return "新增到候选池";
  return "准备导入";
}

function numberText(value, fallback = "-") {
  return Number.isFinite(value) ? String(value) : fallback;
}

function renderResearchPreview(result, imported = false) {
  const research = result.research ?? {};
  const valuation = research.valuation ?? {};
  const quality = research.quality ?? {};
  const warnings = result.warnings ?? [];
  const plan = result.plan ?? [];
  const targetPlan = plan.find((item) => {
    const itemSymbol = String(item.symbol ?? "").toUpperCase();
    return itemSymbol ? itemSymbol === String(research.symbol ?? "").toUpperCase() : item.name === research.name;
  });

  elements.researchPreview.innerHTML = `
    <div class="research-summary">
      <strong>${escapeHTML(targetTypeText(result.targetType))}</strong>
      <span>${escapeHTML(result.summary ?? "")}</span>
      ${result.backupPath ? `<small>备份：${escapeHTML(result.backupPath)}</small>` : ""}
    </div>
    <div class="research-preview-grid">
      <div><span>股票</span><strong>${escapeHTML(research.name ?? "-")}</strong><small>${escapeHTML(research.symbol ?? "-")} · ${escapeHTML(research.asOf ?? "-")}</small></div>
      <div><span>安全边际</span><strong>${Number.isFinite(valuation.marginOfSafety) ? percent(valuation.marginOfSafety * 100, false) : "-"}</strong><small>${escapeHTML(valuation.fairValueRange ?? "-")}</small></div>
      <div><span>质量总分</span><strong>${numberText(quality.totalScore)}</strong><small>${numberText(quality.businessModel)}/${numberText(quality.moat)}/${numberText(quality.governance)}/${numberText(quality.financialQuality)}</small></div>
      <div><span>执行排序</span><strong>${targetPlan ? targetPlan.rank : "-"}</strong><small>${escapeHTML(targetPlan ? targetPlan.priority : "未列入 Plan")}</small></div>
    </div>
    ${warnings.length ? `
      <div class="research-warnings">
        ${warnings.map((item) => `<span>${escapeHTML(item)}</span>`).join("")}
      </div>
    ` : ""}
    <div class="preview-plan">
      ${plan.map((item) => `
        <div class="${item.name === research.name ? "active" : ""}">
          <span>${item.rank}</span>
          <strong>${escapeHTML(item.name)}</strong>
          <small>${escapeHTML(item.priority)}</small>
        </div>
      `).join("")}
    </div>
  `;

  setResearchStatus(imported ? "已导入并刷新页面数据" : "校验通过，可以确认导入", imported ? "success" : "ready");
}

function computePositions() {
  return state.holdings
    .filter((holding) => holding.shares > 0)
    .map((holding) => {
      const marketValueLocal = holding.shares * holding.currentPrice;
      const costValueLocal = holding.shares * holding.cost;
      const marketValueCny = marketValueLocal * fx(holding.currency);
      const costValueCny = costValueLocal * fx(holding.currency);
      const pnlCny = marketValueCny - costValueCny;
      const previousClose = holding.previousClose > 0 ? holding.previousClose : holding.currentPrice;
      const dayChange = holding.shares * (holding.currentPrice - previousClose) * fx(holding.currency);

      return {
        ...holding,
        marginOfSafety: calculatedMarginOfSafety(holding) ?? holding.marginOfSafety,
        marketValueLocal,
        marketValueCny,
        costValueCny,
        pnlCny,
        pnlRate: costValueCny ? (pnlCny / costValueCny) * 100 : 0,
        dayChange
      };
    });
}

function syncCash() {
  if (Number.isFinite(state.cash)) return;

  const invested = state.holdings.reduce((sum, holding) => {
    return sum + holding.shares * holding.currentPrice * fx(holding.currency);
  }, 0);

  state.cash = state.totalCapital - invested;
}

function getFilteredPositions(positions) {
  return positions.filter((position) => {
    const haystack = [
      position.symbol,
      position.name,
      position.action,
      position.status,
      position.industry
    ].join(" ").toLowerCase();
    const matchesSearch = haystack.includes(searchTerm);
    const matchesFilter =
      activeFilter === "all" ||
      (activeFilter === "profit" && position.pnlCny >= 0) ||
      (activeFilter === "loss" && position.pnlCny < 0);

    return matchesSearch && matchesFilter;
  });
}

function holdingHealth(position, totalValue) {
  const margin = calculatedMarginOfSafety(position) ?? finiteNumber(position.marginOfSafety);
  const quality = finiteNumber(position.qualityScore);
  const weight = totalValue ? position.marketValueCny / totalValue : 0;
  const text = [position.risk, position.status, position.action].join(" ");
  const highRisk = /停牌|重大风险|否决|调查|内控|退市|财报可信|风险暴露/.test(text);
  const reduceSignal = /减仓|降权|剔除|风险观察|不纳入核心/.test(text);
  const marginScore = Number.isFinite(margin) ? clamp(((margin + 0.1) / 0.35) * 100, 0, 100) : 45;
  const qualityScore = Number.isFinite(quality) ? quality : 60;
  const pnlScore = clamp(((position.pnlRate ?? 0) + 20) / 40 * 100, 0, 100);
  const weightScore = weight <= 0.15 ? 100 : weight <= 0.25 ? 72 : 45;
  const riskScore = highRisk ? 0 : reduceSignal ? 35 : 100;
  const score = Math.round(qualityScore * 0.35 + marginScore * 0.3 + pnlScore * 0.15 + weightScore * 0.1 + riskScore * 0.1);

  if (highRisk || (Number.isFinite(quality) && quality < 75)) {
    return {
      status: "风险暴露",
      tone: "risk",
      score,
      detail: "风险标签或质量分不达标"
    };
  }
  if (reduceSignal || (Number.isFinite(margin) && margin < 0) || position.pnlRate <= -15) {
    return {
      status: "降权",
      tone: "reduce",
      score,
      detail: "进入降权或复盘区"
    };
  }
  if (Number.isFinite(quality) && quality >= 85 && Number.isFinite(margin) && margin >= 0.15 && weight <= 0.2) {
    return {
      status: "核心",
      tone: "core",
      score,
      detail: `${percent(weight * 100, false)}仓位 · ${displayMarginOfSafety(position)}安全边际`
    };
  }
  return {
    status: "观察",
    tone: "watch",
    score,
    detail: `${percent(weight * 100, false)}仓位 · ${displayMarginOfSafety(position)}安全边际`
  };
}

function renderMetrics(positions) {
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const totalCost = positions.reduce((sum, item) => sum + item.costValueCny, 0);
  const totalPnl = totalValue - totalCost;
  const dayChange = positions.reduce((sum, item) => sum + item.dayChange, 0);

  elements.totalValue.textContent = currency(totalValue);
  elements.totalPnl.textContent = currency(totalPnl);
  elements.totalPnl.className = totalPnl >= 0 ? "positive" : "negative";
  elements.totalPnlRate.textContent = percent(totalCost ? (totalPnl / totalCost) * 100 : 0);
  elements.totalPnlRate.className = totalPnl >= 0 ? "positive" : "negative";
  elements.dayChange.textContent = currency(dayChange);
  elements.dayChange.className = dayChange >= 0 ? "positive" : "negative";
  elements.dayChangeRate.textContent = percent(totalValue ? (dayChange / totalValue) * 100 : 0);
  elements.dayChangeRate.className = dayChange >= 0 ? "positive" : "negative";
  elements.cashBalance.textContent = currency(state.cash);
  elements.positionCount.textContent = `${positions.length} 只股票`;
  elements.recordCount.textContent = `${state.holdings.length} 条持仓 · ${state.trades.length} 条交易`;
}

function renderPositions(positions) {
  const filtered = getFilteredPositions(positions);
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);

  if (!filtered.length) {
    elements.positionsBody.innerHTML = `<tr><td colspan="11" class="empty-state">暂无符合条件的持仓</td></tr>`;
    return;
  }

  elements.positionsBody.innerHTML = filtered
    .map((position) => {
      const weight = totalValue ? (position.marketValueCny / totalValue) * 100 : 0;
      const pnlClass = position.pnlCny >= 0 ? "positive" : "negative";
      const marginText = displayMarginOfSafety(position);
      const qualityText = Number.isFinite(position.qualityScore) ? position.qualityScore : "-";
      const health = holdingHealth(position, totalValue);

      return `
        <tr>
          <td>
            <div class="stock-cell">
              <span class="ticker">${position.symbol.slice(0, 4)}</span>
              <a class="stock-name stock-link" href="${stockHash(position.symbol)}">
                <strong>${escapeHTML(position.name)}</strong>
                <span>${escapeHTML(position.symbol)} · ${escapeHTML(position.industry || "未分类")}</span>
              </a>
            </div>
          </td>
          <td>${position.shares}</td>
          <td>${currency(position.cost, position.currency)}</td>
          <td>
            ${currency(position.currentPrice, position.currency)}
            <br />
            <small class="quote-date">${position.currentPriceDate || "收盘日未知"}</small>
          </td>
          <td>${currency(position.marketValueCny)}</td>
          <td class="${pnlClass}">
            ${currency(position.pnlCny)}
            <br />
            <small>${percent(position.pnlRate)}</small>
          </td>
          <td>${marginText}</td>
          <td>${qualityText}</td>
          <td>
            <div class="weight-bar" title="${weight.toFixed(2)}%">
              <span style="width: ${Math.min(weight, 100)}%"></span>
            </div>
          </td>
          <td>
            <span class="health-pill ${health.tone}" title="${escapeHTML(health.detail)}">${health.status}</span>
            <br />
            <small class="health-score">${health.score} 分</small>
          </td>
          <td>
            <button class="icon-button edit-holding" data-symbol="${position.symbol}" title="编辑 Excel 信息">✎</button>
          </td>
        </tr>
      `;
    })
    .join("");
}

function renderAllocation(positions) {
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);

  elements.allocationChart.innerHTML = positions
    .sort((a, b) => b.marketValueCny - a.marketValueCny)
    .map((position, index) => {
      const share = totalValue ? (position.marketValueCny / totalValue) * 100 : 0;
      return `
        <a class="allocation-row allocation-link" href="${stockHash(position.symbol)}">
          <div class="allocation-label">
            <strong>${escapeHTML(position.name)}</strong>
            <span>${escapeHTML(position.symbol)}</span>
          </div>
          <div class="allocation-track">
            <span style="width: ${share}%; background: ${palette[index % palette.length]}"></span>
          </div>
          <span>${share.toFixed(1)}%</span>
        </a>
      `;
    })
    .join("");
}

function renderTrades() {
  const recentTrades = [...state.trades].reverse();

  if (!recentTrades.length) {
    elements.tradeList.innerHTML = `<div class="empty-state">暂无交易记录</div>`;
    return;
  }

  elements.tradeList.innerHTML = recentTrades
    .map((trade) => {
      const sideClass = trade.side === "buy" ? "positive" : "negative";
      const sideText = trade.side === "buy" ? "买入" : "卖出";
      return `
        <div class="trade-item">
          <strong>${trade.symbol} · ${sideText}</strong>
          <span class="${sideClass}">${currency(trade.price, trade.currency)}</span>
          <small>${trade.date} · ${trade.name}</small>
          <small>${trade.shares} 股 · 最新价 ${currency(trade.currentPrice, trade.currency)}</small>
        </div>
      `;
    })
    .join("");
}

function findStockForPlanItem(item) {
  const symbol = findSymbolForPlan(item);
  const normalized = normalizeSymbol(symbol);
  return (
    state.holdings.find((stock) => normalizeSymbol(stock.symbol) === normalized) ||
    state.candidates.find((stock) => normalizeSymbol(stock.symbol) === normalized) ||
    state.holdings.find((stock) => stock.name === item.name || stock.name.includes(item.name) || item.name.includes(stock.name)) ||
    state.candidates.find((stock) => stock.name === item.name || stock.name.includes(item.name) || item.name.includes(stock.name)) ||
    null
  );
}

function sortedPlanItems() {
  return state.plan
    .map((item) => {
      const stock = findStockForPlanItem(item);
      const margin = stock ? calculatedMarginOfSafety(stock) ?? finiteNumber(stock.marginOfSafety) : null;
      const quality = stock ? finiteNumber(stock.qualityScore) : null;
      return { item, stock, margin, quality };
    })
    .sort((a, b) => {
      const marginA = Number.isFinite(a.margin) ? a.margin : -Infinity;
      const marginB = Number.isFinite(b.margin) ? b.margin : -Infinity;
      return marginB - marginA || (b.quality ?? -Infinity) - (a.quality ?? -Infinity) || a.item.name.localeCompare(b.item.name, "zh-CN");
    });
}

function renderPlanAndCandidates() {
  const plans = sortedPlanItems();
  elements.overviewPlanList.innerHTML = plans.length
    ? plans.map(({ item, stock, margin }, index) => `
      <a class="plan-card" href="${stockHash(stock?.symbol || findSymbolForPlan(item))}">
        <span class="plan-rank">${index + 1}</span>
        <div>
          <strong>${escapeHTML(item.name)}</strong>
          <small>${escapeHTML(item.priority)} · 安全边际 ${Number.isFinite(margin) ? percent(margin * 100, false) : "-"}</small>
        </div>
        <p>${escapeHTML(item.advice)}</p>
      </a>
    `)
    .join("")
    : `<div class="empty-state compact-empty">暂无执行计划</div>`;

  elements.candidateList.innerHTML = sortedCandidates()
    .map((item) => {
      const plan = findPlanForStock(item);
      return `
        <a class="candidate-card" href="${stockHash(item.symbol)}">
          <div class="candidate-head">
            <div>
              <strong>${escapeHTML(item.name)}</strong>
              <span>${escapeHTML(item.symbol)} · ${escapeHTML(item.industry)}</span>
            </div>
            <em>${escapeHTML(plan ? `#${plan.rank} ${plan.priority}` : item.status)}</em>
          </div>
          <div class="candidate-metrics">
            <span>质量 <strong>${Number.isFinite(item.qualityScore) ? item.qualityScore : "-"}</strong></span>
            <span>最新价 <strong>${Number.isFinite(item.currentPrice) && item.currentPrice > 0 ? currency(item.currentPrice, item.currency) : "-"}</strong></span>
            <span>安全边际 <strong>${displayMarginOfSafety(item)}</strong></span>
            <span>买入价 <strong>${Number.isFinite(item.targetBuyPrice) ? currency(item.targetBuyPrice, item.currency) : "-"}</strong></span>
            <span>距买点 <strong>${displayBuyDistance(item)}</strong></span>
          </div>
          <p>${escapeHTML(item.action)}</p>
        </a>
      `;
    })
    .join("");
}

function decisionUniverse(positions) {
  const seen = new Set();
  const holdings = positions.map((position) => {
    const symbol = normalizeSymbol(position.symbol);
    seen.add(symbol);
    return { ...position, sourceType: "holding" };
  });
  const candidates = state.candidates
    .filter((candidate) => {
      const symbol = normalizeSymbol(candidate.symbol);
      return symbol && !seen.has(symbol);
    })
    .map((candidate) => ({
      ...candidate,
      sourceType: "candidate",
      currentPrice: finiteNumber(candidate.currentPrice),
      previousClose: finiteNumber(candidate.previousClose),
      marginOfSafety: calculatedMarginOfSafety(candidate) ?? candidate.marginOfSafety
    }));

  return [...holdings, ...candidates].filter((stock) => stock.symbol && stock.name);
}

function localPrice(stock, key) {
  const value = finiteNumber(stock?.[key]);
  return value && value > 0 ? currency(value, stock.currency ?? "CNY") : "-";
}

function planRank(stock) {
  const plan = findPlanForStock(stock);
  return Number.isFinite(plan?.rank) ? plan.rank : 999;
}

function candidateBuyDistance(candidate) {
  const currentPrice = finiteNumber(candidate.currentPrice);
  const targetBuyPrice = finiteNumber(candidate.targetBuyPrice);
  if (!currentPrice || !targetBuyPrice) return null;
  return (currentPrice - targetBuyPrice) / targetBuyPrice;
}

function displayBuyDistance(candidate) {
  const distance = candidateBuyDistance(candidate);
  if (!Number.isFinite(distance)) return "-";
  return distance <= 0 ? "已到买点" : percent(distance * 100, false);
}

function sortedCandidates() {
  const holdingSymbols = new Set(state.holdings.filter((holding) => holding.shares > 0).map((holding) => normalizeSymbol(holding.symbol)));
  return state.candidates
    .filter((candidate) => !holdingSymbols.has(normalizeSymbol(candidate.symbol)))
    .sort((a, b) => {
      const byName = () => a.name.localeCompare(b.name, "zh-CN");
      if (candidateSort === "margin") {
        const marginA = calculatedMarginOfSafety(a) ?? finiteNumber(a.marginOfSafety) ?? -Infinity;
        const marginB = calculatedMarginOfSafety(b) ?? finiteNumber(b.marginOfSafety) ?? -Infinity;
        return marginB - marginA || planRank(a) - planRank(b) || byName();
      }
      if (candidateSort === "quality") {
        return (finiteNumber(b.qualityScore) ?? -Infinity) - (finiteNumber(a.qualityScore) ?? -Infinity) || planRank(a) - planRank(b) || byName();
      }
      if (candidateSort === "buyDistance") {
        const distanceA = candidateBuyDistance(a);
        const distanceB = candidateBuyDistance(b);
        return (Number.isFinite(distanceA) ? distanceA : Infinity) - (Number.isFinite(distanceB) ? distanceB : Infinity) || planRank(a) - planRank(b) || byName();
      }
      if (candidateSort === "industry") {
        return firstIndustry(a.industry).localeCompare(firstIndustry(b.industry), "zh-CN") || planRank(a) - planRank(b) || byName();
      }
      return planRank(a) - planRank(b) || byName();
    });
}

function buildOpportunitySignals(positions) {
  return decisionUniverse(positions)
    .map((stock) => {
      const currentPrice = finiteNumber(stock.currentPrice);
      const targetBuyPrice = finiteNumber(stock.targetBuyPrice);
      const intrinsicValue = finiteNumber(stock.intrinsicValue);
      const marginOfSafety = calculatedMarginOfSafety(stock) ?? finiteNumber(stock.marginOfSafety);
      const reasons = [];
      let priority = 99;
      let tone = "watch";

      if (currentPrice && targetBuyPrice && currentPrice <= targetBuyPrice) {
        reasons.push("进入买入区");
        priority = Math.min(priority, 1);
        tone = "buy";
      } else if (currentPrice && targetBuyPrice && (currentPrice - targetBuyPrice) / targetBuyPrice <= BUY_PROXIMITY) {
        reasons.push("接近买点");
        priority = Math.min(priority, 2);
      }

      if (Number.isFinite(marginOfSafety) && marginOfSafety >= SAFETY_MARGIN_TARGET) {
        reasons.push("安全边际充足");
        priority = Math.min(priority, 3);
        if (tone !== "buy") tone = "safe";
      }

      if (
        stock.sourceType === "holding" &&
        currentPrice &&
        intrinsicValue &&
        currentPrice >= intrinsicValue &&
        Number.isFinite(marginOfSafety) &&
        marginOfSafety < 0
      ) {
        reasons.push("复盘/减仓");
        priority = Math.min(priority, 4);
        tone = "reduce";
      }

      return { stock, reasons, priority, tone, marginOfSafety };
    })
    .filter((item) => item.reasons.length)
    .sort((a, b) => a.priority - b.priority || planRank(a.stock) - planRank(b.stock) || a.stock.name.localeCompare(b.stock.name, "zh-CN"));
}

function renderDecisionItem(signal) {
  const { stock, reasons, tone, marginOfSafety } = signal;
  const sourceText = stock.sourceType === "holding" ? "持仓" : "候选";

  return `
    <a class="decision-item ${tone}" href="${stockHash(stock.symbol)}">
      <div class="decision-title">
        <strong>${escapeHTML(stock.name)}</strong>
        <span>${escapeHTML(sourceText)} · ${escapeHTML(stock.symbol)}</span>
      </div>
      <div class="decision-tags">
        ${reasons.map((reason) => `<span>${escapeHTML(reason)}</span>`).join("")}
      </div>
      <div class="decision-metrics">
        <span>现价 <strong>${localPrice(stock, "currentPrice")}</strong></span>
        <span>买入价 <strong>${localPrice(stock, "targetBuyPrice")}</strong></span>
        <span>内在价值 <strong>${localPrice(stock, "intrinsicValue")}</strong></span>
        <span>安全边际 <strong>${Number.isFinite(marginOfSafety) ? percent(marginOfSafety * 100, false) : "-"}</strong></span>
      </div>
      <p>${escapeHTML(displayText(stock.action, stock.status))}</p>
    </a>
  `;
}

function renderOpportunityRadar(positions) {
  const signals = buildOpportunitySignals(positions);
  elements.opportunityRadar.innerHTML = signals.length
    ? signals.slice(0, 6).map(renderDecisionItem).join("")
    : `<div class="empty-state compact-empty">暂无进入买点的标的</div>`;
}

function firstIndustry(industry) {
  return String(industry ?? "").split("/").map((item) => item.trim()).find(Boolean) || "未分类";
}

function topGroups(items, keyGetter, valueGetter, limit = 3) {
  const groups = new Map();
  items.forEach((item) => {
    const key = keyGetter(item);
    groups.set(key, (groups.get(key) ?? 0) + valueGetter(item));
  });
  return [...groups.entries()]
    .map(([name, value]) => ({ name, value }))
    .sort((a, b) => b.value - a.value)
    .slice(0, limit);
}

function renderDisciplineDashboard(positions) {
  const investedValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const cashValue = finiteNumber(state.cash) ?? 0;
  const totalAssets = investedValue + cashValue;
  const maxPosition = positions.reduce((max, item) => (item.marketValueCny > (max?.marketValueCny ?? 0) ? item : max), null);
  const lowSafety = positions.filter((item) => {
    const margin = calculatedMarginOfSafety(item) ?? finiteNumber(item.marginOfSafety);
    return !Number.isFinite(margin) || margin < SAFETY_MARGIN_TARGET;
  });
  const lowSafetyValue = lowSafety.reduce((sum, item) => sum + item.marketValueCny, 0);
  const industryText = topGroups(positions, (item) => firstIndustry(item.industry), (item) => item.marketValueCny)
    .map((item) => `${item.name} ${percent(investedValue ? (item.value / investedValue) * 100 : 0, false)}`)
    .join(" · ") || "-";
  const currencyText = topGroups(positions, (item) => item.currency ?? "CNY", (item) => item.marketValueCny, 4)
    .map((item) => `${item.name} ${percent(investedValue ? (item.value / investedValue) * 100 : 0, false)}`)
    .join(" · ") || "-";

  const items = [
    {
      label: "现金比例",
      value: percent(totalAssets ? (cashValue / totalAssets) * 100 : 0, false),
      detail: `现金 ${currency(cashValue)}`
    },
    {
      label: "单股最大仓位",
      value: maxPosition ? percent(totalAssets ? (maxPosition.marketValueCny / totalAssets) * 100 : 0, false) : "-",
      detail: maxPosition ? `${maxPosition.name} · ${currency(maxPosition.marketValueCny)}` : "-"
    },
    {
      label: "行业集中度 Top 3",
      value: industryText,
      detail: "按持仓市值计算"
    },
    {
      label: "币种暴露",
      value: currencyText,
      detail: "按持仓市值折人民币计算"
    },
    {
      label: "安全边际不足",
      value: `${lowSafety.length}/${positions.length} 只`,
      detail: `占持仓市值 ${percent(investedValue ? (lowSafetyValue / investedValue) * 100 : 0, false)}`
    }
  ];

  elements.disciplineDashboard.innerHTML = items
    .map((item) => `
      <div class="discipline-item">
        <span>${escapeHTML(item.label)}</span>
        <strong>${escapeHTML(item.value)}</strong>
        <small>${escapeHTML(item.detail)}</small>
      </div>
    `)
    .join("");
}

function renderTriggerAlerts(positions) {
  const signals = buildOpportunitySignals(positions);
  elements.triggerAlerts.innerHTML = signals.length
    ? signals.slice(0, 5).map((signal) => {
        const stock = signal.stock;
        return `
          <a class="trigger-item ${signal.tone}" href="${stockHash(stock.symbol)}">
            <strong>${escapeHTML(stock.name)}</strong>
            <span>${signal.reasons.map(escapeHTML).join(" / ")}</span>
            <small>现价 ${localPrice(stock, "currentPrice")} · 买入价 ${localPrice(stock, "targetBuyPrice")} · 内在价值 ${localPrice(stock, "intrinsicValue")} · 安全边际 ${Number.isFinite(signal.marginOfSafety) ? percent(signal.marginOfSafety * 100, false) : "-"}</small>
            <p>${escapeHTML(displayText(stock.action, stock.status))}</p>
          </a>
        `;
      }).join("")
    : `<div class="empty-state compact-empty">暂无行情触发提醒</div>`;
}

function decisionLogTypeText(type) {
  if (type === "research") return "导入分析";
  if (type === "quote") return "更新行情";
  if (type === "trade") return "新增交易";
  return "记录";
}

function decisionLogTone(type) {
  if (type === "research") return "research";
  if (type === "quote") return "quote";
  if (type === "trade") return "trade";
  return "";
}

function sortedDecisionLogs(symbol = "") {
  const normalizedSymbol = String(symbol ?? "").toUpperCase();
  return [...(state.decisionLogs ?? [])]
    .filter((log) => !normalizedSymbol || String(log.symbol ?? "").toUpperCase() === normalizedSymbol)
    .sort((a, b) => String(b.date ?? "").localeCompare(String(a.date ?? "")) || (Number(b.id) || 0) - (Number(a.id) || 0));
}

function renderDecisionLogItems(logs, emptyText) {
  if (!logs.length) {
    return `<div class="empty-state compact-empty">${escapeHTML(emptyText)}</div>`;
  }

  return logs
    .map((log) => {
      const price = finiteNumber(log.price);
      const priceText = Number.isFinite(price) ? currency(price, log.currency || "CNY") : "价格未记录";
      const title = displayText(log.name, log.symbol || "未知标的");
      return `
        <div class="timeline-item ${decisionLogTone(log.type)}">
          <div class="timeline-head">
            <strong>${escapeHTML(title)}</strong>
            <span>${escapeHTML(displayText(log.date, "时间未知"))}</span>
          </div>
          <div class="timeline-meta">
            <span>${escapeHTML(decisionLogTypeText(log.type))}</span>
            <span>${escapeHTML(displayText(log.symbol, "-"))}</span>
            <span>${priceText}</span>
          </div>
          <p>${escapeHTML(displayText(log.decision, "未记录判断"))}</p>
          <small>${escapeHTML(displayText(log.discipline, "未记录纪律"))}</small>
          ${log.detail ? `<em>${escapeHTML(log.detail)}</em>` : ""}
        </div>
      `;
    })
    .join("");
}

function renderDecisionLogs() {
  elements.decisionLogList.innerHTML = renderDecisionLogItems(sortedDecisionLogs().slice(0, 10), "暂无决策日志");
}

function renderStockDecisionLogs(stock) {
  return renderDecisionLogItems(sortedDecisionLogs(stock.symbol).slice(0, 8), "暂无该股票的决策日志");
}

function renderDecisionArea(positions) {
  renderOpportunityRadar(positions);
  renderDisciplineDashboard(positions);
  renderTriggerAlerts(positions);
}

function findSymbolForPlan(item) {
  const symbol = String(item.symbol ?? "").trim();
  return symbol ? normalizeSymbol(symbol) : findSymbolByName(item.name);
}

function findSymbolByName(name) {
  const normalized = String(name ?? "").trim();
  const holding = state.holdings.find((item) => item.name === normalized || item.name.includes(normalized) || normalized.includes(item.name));
  if (holding) return normalizeSymbol(holding.symbol);
  const candidate = state.candidates.find((item) => item.name === normalized || item.name.includes(normalized) || normalized.includes(item.name));
  return candidate ? normalizeSymbol(candidate.symbol) : "";
}

function metricCard(label, value, detail = "") {
  return `
    <article class="metric">
      <span>${escapeHTML(label)}</span>
      <strong>${value}</strong>
      <small>${escapeHTML(detail)}</small>
    </article>
  `;
}

function scoreItem(label, value, maxScore) {
  const text = Number.isFinite(value) ? `${value}/${maxScore}` : "-";
  const width = Number.isFinite(value) && maxScore ? Math.max(0, Math.min((value / maxScore) * 100, 100)) : 0;
  return `
    <div class="score-item">
      <div>
        <strong>${escapeHTML(label)}</strong>
        <span>${text}</span>
      </div>
      <div class="weight-bar"><span style="width: ${width}%"></span></div>
    </div>
  `;
}

function findStockRecord(symbol, positions) {
  const normalized = normalizeSymbol(symbol);
  const position = positions.find((item) => normalizeSymbol(item.symbol) === normalized);
  if (position) return { stock: position, isHolding: true };

  const candidate = state.candidates.find((item) => normalizeSymbol(item.symbol) === normalized);
  if (!candidate) return { stock: null, isHolding: false };

  return {
    stock: {
      ...candidate,
      shares: 0,
      cost: 0,
      currentPrice: candidate.currentPrice ?? 0,
      previousClose: candidate.previousClose ?? 0,
      currentPriceDate: candidate.currentPriceDate ?? "",
      previousCloseDate: candidate.previousCloseDate ?? "",
      marketValueCny: 0,
      pnlCny: 0,
      pnlRate: 0,
      dayChange: 0
    },
    isHolding: false
  };
}

function findPlanForStock(stock) {
  const stockSymbol = normalizeSymbol(stock.symbol);
  return state.plan.find((item) => {
    const itemSymbol = normalizeSymbol(item.symbol);
    if (itemSymbol && stockSymbol) return itemSymbol === stockSymbol;
    return item.name === stock.name || stock.name.includes(item.name) || item.name.includes(stock.name);
  });
}

function renderReportLibrary(stock) {
  const reports = [...(stock.reports ?? [])].sort((a, b) => {
    const dateCompare = String(b.date ?? "").localeCompare(String(a.date ?? ""));
    return dateCompare || String(b.period ?? "").localeCompare(String(a.period ?? ""));
  });

  if (!reports.length) {
    return `<div class="empty-state compact-empty">暂无近两年定期报告 PDF</div>`;
  }

  return `
    <div class="report-list">
      ${reports.map((report) => `
        <a class="report-card" href="${escapeHTML(report.url)}" target="_blank" rel="noopener noreferrer">
          <span>${escapeHTML(displayText(report.period, "报告期未知"))}</span>
          <strong>${escapeHTML(displayText(report.title, "未命名报告"))}</strong>
          <small>${escapeHTML(displayText(report.kind, "财报"))} · ${escapeHTML(displayText(report.date, "日期未知"))} · ${escapeHTML(displayText(report.source, "官方来源"))}</small>
        </a>
      `).join("")}
    </div>
  `;
}

function renderStockDetail(positions, symbol) {
  const { stock, isHolding } = findStockRecord(symbol, positions);

  if (!stock) {
    elements.stockDetail.innerHTML = `
      <section class="panel">
        <div class="empty-state">未找到该股票</div>
      </section>
    `;
    return;
  }

  const plan = findPlanForStock(stock);
  const pnlClass = stock.pnlCny >= 0 ? "positive" : "negative";
  const dayClass = stock.dayChange >= 0 ? "positive" : "negative";
  const marginText = displayMarginOfSafety(stock);
  const qualityText = Number.isFinite(stock.qualityScore) ? `${stock.qualityScore}` : "-";
  const hasCurrentQuote = Number.isFinite(stock.currentPrice) && stock.currentPrice > 0;
  const hasPreviousQuote = Number.isFinite(stock.previousClose) && stock.previousClose > 0;
  const priceChange = hasCurrentQuote && hasPreviousQuote ? stock.currentPrice - stock.previousClose : null;
  const priceChangeClass = Number.isFinite(priceChange) && priceChange >= 0 ? "positive" : "negative";
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const health = isHolding ? holdingHealth(stock, totalValue) : null;

  elements.stockDetail.innerHTML = `
    <section class="detail-hero">
      <a class="ghost-button detail-back" href="#positions">返回持仓</a>
      <div>
        <p class="eyebrow">${escapeHTML(stock.symbol)} · ${escapeHTML(displayText(stock.industry, "未分类"))}</p>
        <h2>${escapeHTML(stock.name)}</h2>
      </div>
      <div class="detail-hero-meta">
        <span>${escapeHTML(closeDateText(stock) || displayText(stock.updatedAt, "行情日期未知"))}</span>
      </div>
    </section>

    <section class="metrics-grid detail-metrics">
      ${metricCard("今天收盘", hasCurrentQuote ? currency(stock.currentPrice, stock.currency) : "-", stock.currentPriceDate || "收盘日未知")}
      ${metricCard("昨天收盘", hasPreviousQuote ? currency(stock.previousClose, stock.currency) : "-", stock.previousCloseDate || "收盘日未知")}
      ${metricCard("浮动盈亏", `<span class="${pnlClass}">${isHolding ? currency(stock.pnlCny) : "-"}</span>`, isHolding ? percent(stock.pnlRate) : "")}
      ${metricCard(
        "今日变动",
        isHolding
          ? `<span class="${dayClass}">${currency(stock.dayChange)}</span>`
          : `<span class="${priceChangeClass}">${Number.isFinite(priceChange) ? currency(priceChange, stock.currency) : "-"}</span>`,
        isHolding && stock.marketValueCny
          ? percent((stock.dayChange / stock.marketValueCny) * 100)
          : Number.isFinite(priceChange) && stock.previousClose
            ? percent((priceChange / stock.previousClose) * 100)
            : ""
      )}
    </section>

    <section class="detail-grid">
      <section class="panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Analysis</p>
            <h2>股票分析</h2>
          </div>
        </div>
        <div class="detail-content">
          <dl class="detail-list">
            <div><dt>达标状态</dt><dd>${escapeHTML(displayText(stock.status, "未填写"))}</dd></div>
            <div><dt>最终动作</dt><dd>${escapeHTML(displayText(stock.action, "未填写"))}</dd></div>
            <div><dt>主要风险</dt><dd>${escapeHTML(displayText(stock.risk, "未填写"))}</dd></div>
            <div><dt>备注</dt><dd>${escapeHTML(displayText(stock.notes, "未填写"))}</dd></div>
          </dl>
        </div>
      </section>

      <section class="panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Valuation</p>
            <h2>估值与质量</h2>
          </div>
        </div>
        <div class="detail-content">
          <div class="valuation-grid">
            <div><span>安全边际</span><strong>${marginText}</strong></div>
            <div><span>质量总分</span><strong>${qualityText}</strong></div>
            <div><span>健康状态</span><strong>${health ? `<span class="health-pill ${health.tone}">${health.status}</span>` : "-"}</strong></div>
            <div><span>健康评分</span><strong>${health ? `${health.score} 分` : "-"}</strong></div>
            <div><span>内在价值</span><strong>${Number.isFinite(stock.intrinsicValue) ? currency(stock.intrinsicValue, stock.currency) : "-"}</strong></div>
            <div><span>买入价</span><strong>${Number.isFinite(stock.targetBuyPrice) ? currency(stock.targetBuyPrice, stock.currency) : "-"}</strong></div>
            <div><span>公允区间</span><strong>${escapeHTML(displayText(stock.fairValueRange))}</strong></div>
            <div><span>持仓市值</span><strong>${isHolding ? currency(stock.marketValueCny) : "-"}</strong></div>
          </div>
          <div class="score-list">
            ${scoreItem("商业模式", stock.businessModel, 30)}
            ${scoreItem("护城河", stock.moat, 25)}
            ${scoreItem("治理", stock.governance, 20)}
            ${scoreItem("财务质量", stock.financialQuality, 25)}
          </div>
        </div>
      </section>
    </section>

    <section class="panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Execution</p>
          <h2>执行计划</h2>
        </div>
      </div>
      <div class="detail-content">
        <div class="execution-plan">
          <strong>${plan ? `${plan.rank}. ${escapeHTML(plan.name)} · ${escapeHTML(plan.priority)}` : "未单独列入执行优先级"}</strong>
          <span>${escapeHTML(plan ? plan.advice : displayText(stock.action, "暂无执行计划"))}</span>
          <small>${escapeHTML(plan ? plan.discipline : displayText(stock.status, "暂无纪律说明"))}</small>
        </div>
      </div>
    </section>

    <section class="panel decision-log-panel detail-log-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Timeline</p>
          <h2>决策日志</h2>
        </div>
      </div>
      <div class="decision-log-list">
        ${renderStockDecisionLogs(stock)}
      </div>
    </section>

    <section class="panel report-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Reports</p>
          <h2>近两年财报 PDF</h2>
        </div>
      </div>
      <div class="detail-content">
        ${renderReportLibrary(stock)}
      </div>
    </section>
  `;
}

function render() {
  const positions = computePositions();
  renderMetrics(positions);
  renderPositions(positions);
  renderDecisionArea(positions);
  renderDecisionLogs();
  renderAllocation(positions);
  renderTrades();
  renderPlanAndCandidates();
  if (window.location.hash.startsWith("#stock=")) {
    renderStockDetail(positions, decodeURIComponent(window.location.hash.slice("#stock=".length)));
  }
}

async function addTrade(formData) {
  const side = formData.get("side");
  const shares = Number(formData.get("shares"));
  const price = Number(formData.get("price"));
  const symbol = String(formData.get("symbol")).trim().toUpperCase();
  const currencyCode = String(formData.get("currency"));
  const trade = {
    id: Date.now(),
    date: new Date().toISOString().slice(0, 10),
    symbol,
    name: String(formData.get("name")).trim(),
    side,
    shares,
    price,
    currency: currencyCode,
    currentPrice: Number(formData.get("currentPrice"))
  };

  if (USE_BACKEND) {
    state = await requestJSON("/api/trades", {
      method: "POST",
      body: JSON.stringify(trade)
    });
    saveState();
    render();
    return;
  }

  let holding = state.holdings.find((item) => item.symbol.toUpperCase() === symbol);
  if (!holding) {
    holding = {
      symbol,
      name: trade.name,
      shares: 0,
      cost: price,
      currentPrice: trade.currentPrice,
      previousClose: trade.currentPrice,
      currentPriceDate: new Date().toISOString().slice(0, 10),
      previousCloseDate: new Date().toISOString().slice(0, 10),
      currency: currencyCode,
      action: "",
      status: "",
      marginOfSafety: null,
      qualityScore: null,
      industry: "",
      notes: ""
    };
    state.holdings.push(holding);
  }

  if (side === "buy") {
    const totalCost = holding.shares * holding.cost + shares * price;
    holding.shares += shares;
    holding.cost = totalCost / holding.shares;
  } else {
    holding.shares = Math.max(0, holding.shares - shares);
  }

  holding.name = trade.name || holding.name;
  holding.currency = currencyCode;
  holding.currentPrice = trade.currentPrice;
  holding.previousClose = holding.previousClose > 0 ? holding.previousClose : trade.currentPrice;
  holding.currentPriceDate = holding.currentPriceDate || new Date().toISOString().slice(0, 10);
  holding.previousCloseDate = holding.previousCloseDate || holding.currentPriceDate;
  state.trades.push(trade);
  state.cash += side === "sell" ? shares * price * fx(currencyCode) : -(shares * price * fx(currencyCode));
  const plan = findPlanForStock(holding);
  const sideText = side === "buy" ? "买入" : "卖出";
  appendClientDecisionLog({
    type: "trade",
    symbol,
    name: holding.name,
    price,
    currency: currencyCode,
    decision: `${sideText} ${holding.name}`,
    discipline: plan?.discipline || holding.status || "未记录纪律",
    detail: `${sideText} ${shares} 股；成交价 ${currencyCode} ${price.toFixed(4)}；录入最新价 ${currencyCode} ${trade.currentPrice.toFixed(4)}`
  });
  saveState();
  render();
}

async function previewResearch() {
  if (!USE_BACKEND) throw new Error("需要通过 go run . 启动后端后才能导入到 portfolio.json");

  const research = parseResearchJSON();
  const result = await requestJSON("/api/research/preview", {
    method: "POST",
    body: JSON.stringify(research)
  });
  pendingResearch = research;
  elements.importResearchButton.disabled = false;
  renderResearchPreview(result);
}

async function importResearch() {
  if (!pendingResearch) {
    await previewResearch();
    if (!pendingResearch) return;
  }

  const result = await requestJSON("/api/research/import", {
    method: "POST",
    body: JSON.stringify(pendingResearch)
  });

  state = result.state;
  localStorage.removeItem(STORAGE_KEY);
  pendingResearch = null;
  elements.importResearchButton.disabled = true;
  renderResearchPreview(result, true);
  render();

  if (result.research?.symbol) {
    window.location.hash = stockHash(result.research.symbol);
  }
}

async function updateQuotes() {
  if (!USE_BACKEND) throw new Error("需要通过 go run . 启动后端后才能更新行情");

  elements.updateQuotesButton.disabled = true;
  elements.updateQuotesButton.innerHTML = "<span>↻</span> 更新中";
  setQuoteUpdateStatus("正在拉取持仓和候选池最新日线收盘价...");

  try {
    const result = await requestJSON("/api/quotes/update", { method: "POST" });
    state = result.state;
    localStorage.removeItem(STORAGE_KEY);
    syncCash();
    render();

    const alertCount = buildOpportunitySignals(computePositions()).length;
    const skipped = result.skipped ?? [];
    if (skipped.length) {
      const preview = skipped.slice(0, 3).map((item) => `${item.symbol} ${item.error}`).join("；");
      setQuoteUpdateStatus(`已更新 ${result.updated} 个标的，触发 ${alertCount} 条提醒，${skipped.length} 个失败：${preview}`, "error");
    } else {
      setQuoteUpdateStatus(`已更新 ${result.updated} 个标的，触发 ${alertCount} 条提醒`, "success");
    }
  } finally {
    elements.updateQuotesButton.disabled = false;
    elements.updateQuotesButton.innerHTML = "<span>↻</span> 更新行情";
  }
}

function openHoldingEditor(symbol) {
  const holding = state.holdings.find((item) => item.symbol === symbol);
  if (!holding) return;

  const form = elements.holdingForm;
  form.symbol.value = holding.symbol;
  form.name.value = holding.name ?? "";
  form.industry.value = holding.industry ?? "";
  form.action.value = holding.action ?? "";
  form.status.value = holding.status ?? "";
  form.marginOfSafety.value = Number.isFinite(calculatedMarginOfSafety(holding)) ? calculatedMarginOfSafety(holding) : "";
  form.qualityScore.value = Number.isFinite(holding.qualityScore) ? holding.qualityScore : "";
  form.notes.value = holding.notes ?? "";
  elements.holdingDialog.showModal();
}

async function saveHolding(formData) {
  const symbol = formData.get("symbol");
  const holding = state.holdings.find((item) => item.symbol === symbol);
  if (!holding) return;

  const patch = {
    name: String(formData.get("name")).trim(),
    industry: String(formData.get("industry")).trim(),
    action: String(formData.get("action")).trim(),
    status: String(formData.get("status")).trim(),
    marginOfSafety: calculatedMarginOfSafety(holding),
    qualityScore: formData.get("qualityScore") === "" ? null : Number(formData.get("qualityScore")),
    notes: String(formData.get("notes")).trim()
  };

  if (USE_BACKEND) {
    state = await requestJSON(`/api/holdings/${encodeURIComponent(symbol)}`, {
      method: "PUT",
      body: JSON.stringify(patch)
    });
    saveState();
    render();
    return;
  }

  holding.name = patch.name;
  holding.industry = patch.industry;
  holding.action = patch.action;
  holding.status = patch.status;
  holding.marginOfSafety = patch.marginOfSafety;
  holding.qualityScore = patch.qualityScore;
  holding.notes = patch.notes;
  saveState();
  render();
}

document.querySelectorAll(".segment").forEach((button) => {
  button.addEventListener("click", () => {
    document.querySelector(".segment.active").classList.remove("active");
    button.classList.add("active");
    activeFilter = button.dataset.filter;
    render();
  });
});

document.querySelectorAll(".nav-item").forEach((button) => {
  button.addEventListener("click", () => {
    window.location.hash = button.dataset.view;
  });
});

function showPage(view) {
  const isStockDetail = view.startsWith("stock=");
  const nextView = isStockDetail ? "stock-detail" : pageTitles[view] ? view : "overview";

  if (isStockDetail) {
    renderStockDetail(computePositions(), decodeURIComponent(view.slice("stock=".length)));
  }

  document.querySelector(".nav-item.active")?.classList.remove("active");
  document.querySelector(`.nav-item[data-view="${nextView}"]`)?.classList.add("active");
  document.querySelector(".page.active")?.classList.remove("active");
  document.querySelector(`[data-page="${nextView}"]`)?.classList.add("active");
  elements.pageTitle.textContent = pageTitles[nextView];
}

window.addEventListener("hashchange", () => {
  showPage(window.location.hash.slice(1));
});

elements.searchInput.addEventListener("input", (event) => {
  searchTerm = event.target.value.trim().toLowerCase();
  render();
});

elements.candidateSort.addEventListener("change", (event) => {
  candidateSort = event.target.value;
  renderPlanAndCandidates();
});

document.querySelector("#openTradePanelSecondary").addEventListener("click", () => {
  elements.tradeDialog.showModal();
});

document.querySelector("#openResearchPanel").addEventListener("click", () => {
  pendingResearch = null;
  elements.importResearchButton.disabled = true;
  elements.researchPreview.innerHTML = "";
  setResearchStatus("");
  elements.researchDialog.showModal();
});

elements.updateQuotesButton.addEventListener("click", async () => {
  try {
    await updateQuotes();
  } catch (error) {
    setQuoteUpdateStatus(error.message, "error");
  }
});

document.querySelector("#closeTradePanel").addEventListener("click", () => {
  elements.tradeDialog.close();
});

document.querySelector("#cancelTrade").addEventListener("click", () => {
  elements.tradeDialog.close();
});

document.querySelector("#closeHoldingPanel").addEventListener("click", () => {
  elements.holdingDialog.close();
});

document.querySelector("#cancelHolding").addEventListener("click", () => {
  elements.holdingDialog.close();
});

document.querySelector("#closeResearchPanel").addEventListener("click", () => {
  elements.researchDialog.close();
});

document.querySelector("#cancelResearch").addEventListener("click", () => {
  elements.researchDialog.close();
});

elements.tradeForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  await addTrade(new FormData(elements.tradeForm));
  elements.tradeForm.reset();
  elements.tradeDialog.close();
});

elements.holdingForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  await saveHolding(new FormData(elements.holdingForm));
  elements.holdingDialog.close();
});

elements.researchForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    elements.importResearchButton.disabled = true;
    setResearchStatus("正在校验...");
    await previewResearch();
  } catch (error) {
    pendingResearch = null;
    elements.researchPreview.innerHTML = "";
    setResearchStatus(error.message, "error");
  }
});

elements.importResearchButton.addEventListener("click", async () => {
  try {
    elements.importResearchButton.disabled = true;
    setResearchStatus("正在写入 portfolio.json...");
    await importResearch();
  } catch (error) {
    elements.importResearchButton.disabled = false;
    setResearchStatus(error.message, "error");
  }
});

elements.researchJSON.addEventListener("input", () => {
  pendingResearch = null;
  elements.importResearchButton.disabled = true;
  elements.researchPreview.innerHTML = "";
  setResearchStatus("");
});

elements.positionsBody.addEventListener("click", (event) => {
  const editButton = event.target.closest(".edit-holding");
  if (!editButton) return;
  openHoldingEditor(editButton.dataset.symbol);
});

async function init() {
  await loadBackendState();
  syncCash();
  render();
  showPage(window.location.hash.slice(1));
}

init();
