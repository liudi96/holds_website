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
      name: "海尔智家",
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
const MAIN_DCF_MARGIN_TARGET = 0.15;
const MAIN_ALLOCATION_TARGET = 0.7;
const CIGAR_ALLOCATION_TARGET = 0.3;
const A_SHARE_SHAREHOLDER_RETURN_TARGET = 0.06;
const HK_SHARE_SHAREHOLDER_RETURN_TARGET = 0.08;
const A_SHARE_EX_CASH_PE_MAX = 10;
const HK_SHARE_EX_CASH_PE_MAX = 8;
const AGGRESSIVE_BUY_DISCOUNT = 0.1;
const OWNER_AUDIT_SCORE_TARGET = 75;
const MAJOR_RISK_PATTERN = /停牌|重大风险|否决|调查|内控|退市|财报可信|风险暴露|治理风险|治理与财务可靠性|质量分<75|低于75/;
const OWNER_AUDIT_FIELDS = [
  { key: "tenYearDemand", label: "十年需求", core: true, weight: 18, evidence: "十年后需求是否仍稳定，是否受替代品或政策永久削弱" },
  { key: "assetDurability", label: "资产耐久", core: false, weight: 14, evidence: "品牌、网络、牌照、资源或渠道资产是否能长期保持价值" },
  { key: "maintenanceCapexLight", label: "轻再投资", core: false, weight: 12, evidence: "维持竞争力所需资本开支是否轻，扩张是否吞噬自由现金流" },
  { key: "dividendFcfSupport", label: "分红FCF", core: true, weight: 18, evidence: "分红和回购是否由真实自由现金流覆盖，而非靠借债或卖资产" },
  { key: "dividendReinvestmentEfficiency", label: "再投资效率", core: false, weight: 12, evidence: "留存利润、分红和回购哪种资本配置对股东更有效" },
  { key: "roeRoicDurability", label: "ROE/ROIC", core: false, weight: 14, evidence: "ROE/ROIC 是否可持续，是否依赖高杠杆或周期高点" },
  { key: "valuationSystemRisk", label: "估值体系", core: true, weight: 12, evidence: "行业估值锚是否发生永久变化，当前估值假设是否仍成立" }
];
const OWNER_AUDIT_STATUS = {
  pass: { text: "通过", tone: "strong" },
  review: { text: "复核", tone: "watch" },
  fail: { text: "失败", tone: "risk" }
};
const OWNER_AUDIT_STATUS_SCORE = {
  pass: 1,
  review: 0.6,
  fail: 0
};

let state = loadState();
let activeFilter = "all";
let searchTerm = "";
let pendingResearch = null;
let candidateSort = "consensus";
let candidateFilter = "all";
let decisionLogFilter = "all";
let masterMatrixSort = { key: "margin", direction: "desc" };
const pageTitles = {
  overview: "总览",
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
  committeeConsensus: document.querySelector("#committeeConsensus"),
  opportunityRadar: document.querySelector("#opportunityRadar"),
  disciplineDashboard: document.querySelector("#disciplineDashboard"),
  dataQualityList: document.querySelector("#dataQualityList"),
  dividendCashList: document.querySelector("#dividendCashList"),
  decisionLogList: document.querySelector("#decisionLogList"),
  decisionLogPanel: document.querySelector("#decisionLogPanel"),
  decisionLogToggle: document.querySelector("#decisionLogToggle"),
  decisionLogFilters: document.querySelector("#decisionLogFilters"),
  clearDecisionLogs: document.querySelector("#clearDecisionLogs"),
  masterMatrix: document.querySelector("#masterMatrix"),
  grahamSummary: document.querySelector("#grahamSummary"),
  grahamList: document.querySelector("#grahamList"),
  buffettSummary: document.querySelector("#buffettSummary"),
  buffettList: document.querySelector("#buffettList"),
  candidateList: document.querySelector("#candidateList"),
  candidateSort: document.querySelector("#candidateSort"),
  stockDetail: document.querySelector("#stockDetail"),
  totalFunds: document.querySelector("#totalFunds"),
  totalValue: document.querySelector("#totalValue"),
  dayChange: document.querySelector("#dayChange"),
  dayChangeRate: document.querySelector("#dayChangeRate"),
  annualDividend: document.querySelector("#annualDividend"),
  portfolioDividendYield: document.querySelector("#portfolioDividendYield"),
  dataQualityMetric: document.querySelector("#dataQualityMetric"),
  dataQualityDetail: document.querySelector("#dataQualityDetail"),
  actionConclusion: document.querySelector("#actionConclusion"),
  actionConclusionStatus: document.querySelector("#actionConclusionStatus"),
  actionConclusionDetail: document.querySelector("#actionConclusionDetail"),
  positionCount: document.querySelector("#positionCount"),
  recordCount: document.querySelector("#recordCount"),
  searchInput: document.querySelector("#searchInput"),
  searchResults: document.querySelector("#searchResults"),
  updateQuotesButton: document.querySelector("#updateQuotesButton"),
  exportChatGPTButton: document.querySelector("#exportChatGPTContext"),
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

async function clearNonTradeDecisionLogs() {
  const confirmed = window.confirm("将清理分析、行情等非交易决策日志，并保留新增交易日志。此操作会写回 portfolio.json，是否继续？");
  if (!confirmed) return;

  if (USE_BACKEND) {
    state = await requestJSON("/api/decision-logs/clear", { method: "POST" });
  } else {
    state.decisionLogs = [...(state.decisionLogs ?? [])].filter((log) => String(log.type ?? "").toLowerCase() === "trade");
    saveState();
  }

  render();
}

function fx(currencyCode) {
  return state.fx[currencyCode] ?? 1;
}

function currency(value, currencyCode = "CNY") {
  const number = finiteNumber(value);
  if (number === null) return "-";
  return new Intl.NumberFormat("zh-CN", {
    style: "currency",
    currency: currencyCode,
    minimumFractionDigits: 2
  }).format(number);
}

function wholeCurrency(value, currencyCode = "CNY") {
  const number = finiteNumber(value);
  if (number === null) return "-";
  return new Intl.NumberFormat("zh-CN", {
    style: "currency",
    currency: currencyCode,
    maximumFractionDigits: 0
  }).format(number);
}

function percent(value, signed = true) {
  const number = finiteNumber(value);
  if (number === null) return "-";
  const prefix = signed && number >= 0 ? "+" : "";
  return `${prefix}${number.toFixed(2)}%`;
}

function finiteNumber(value) {
  if (value === null || value === undefined || value === "") return null;
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

function computedInitialBuyPrice(stock) {
  const intrinsicValue = finiteNumber(stock?.intrinsicValue);
  if (Number.isFinite(intrinsicValue) && intrinsicValue > 0) {
    return intrinsicValue * (1 - SAFETY_MARGIN_TARGET);
  }
  const targetBuyPrice = finiteNumber(stock?.targetBuyPrice);
  if (Number.isFinite(targetBuyPrice) && targetBuyPrice > 0) return targetBuyPrice;
  return null;
}

function priceLevels(stock) {
  const initialBuyPrice = computedInitialBuyPrice(stock);
  if (!Number.isFinite(initialBuyPrice) || initialBuyPrice <= 0) {
    return { watchPrice: null, initialBuyPrice: null, aggressiveBuyPrice: null };
  }
  return {
    watchPrice: initialBuyPrice * (1 + BUY_PROXIMITY),
    initialBuyPrice,
    aggressiveBuyPrice: initialBuyPrice * (1 - AGGRESSIVE_BUY_DISCOUNT)
  };
}

function displayPriceLevel(stock, key) {
  const value = priceLevels(stock)[key];
  return Number.isFinite(value) ? currency(value, stock.currency ?? "CNY") : "-";
}

function riskText(stock) {
  return [stock?.risk, stock?.status, stock?.action, stock?.notes].filter(Boolean).join(" ");
}

function hasMajorRisk(stock) {
  const text = riskText(stock).replace(/无一票否决|无立即否决|没有一票否决/g, "");
  return MAJOR_RISK_PATTERN.test(text);
}

function valuationConfidence(stock) {
  const explicit = String(stock?.valuationConfidence ?? "").trim().toLowerCase();
  if (["high", "medium", "low"].includes(explicit)) return explicit;
  const quality = finiteNumber(stock?.qualityScore);
  if (hasMajorRisk(stock) || (Number.isFinite(quality) && quality < 75)) return "low";
  if (Number.isFinite(quality) && quality >= 85) return "high";
  return "medium";
}

function confidenceMeta(stock) {
  const value = valuationConfidence(stock);
  if (value === "high") return { value, text: "高可信", tone: "strong" };
  if (value === "low") return { value, text: "低可信", tone: "risk" };
  return { value, text: "中可信", tone: "watch" };
}

function confidenceScore(stock) {
  const confidence = valuationConfidence(stock);
  if (confidence === "high") return 1;
  if (confidence === "medium") return 0.7;
  return 0.2;
}

function badge(text, tone = "watch") {
  return `<span class="status-badge ${tone}">${escapeHTML(text)}</span>`;
}

function dividendCurrency(stock) {
  return String(stock?.dividend?.dividendCurrency || stock?.currency || "CNY").toUpperCase();
}

function cashDividendTotalCurrency(stock) {
  return String(stock?.dividend?.cashDividendCurrency || stock?.dividend?.dividendCurrency || stock?.currency || "CNY").toUpperCase();
}

function buybackCurrency(stock) {
  return String(stock?.dividend?.buybackCurrency || stock?.currency || "CNY").toUpperCase();
}

function marketCapCurrency(stock) {
  return String(stock?.marketCapCurrency || stock?.currency || "CNY").toUpperCase();
}

function dividendYieldInputs(stock) {
  const dividend = stock?.dividend;
  const cashDividendTotal = finiteNumber(dividend?.cashDividendTotal);
  const marketCap = finiteNumber(stock?.marketCap);
  if (Number.isFinite(cashDividendTotal) && cashDividendTotal > 0 && Number.isFinite(marketCap) && marketCap > 0) {
    return {
      cashDividendTotalCny: cashDividendTotal * fx(cashDividendTotalCurrency(stock)),
      marketCapCny: marketCap * fx(marketCapCurrency(stock))
    };
  }

  // Same-share fallback: total cash dividend / total market cap reduces to DPS / price.
  const perShare = finiteNumber(dividend?.dividendPerShare);
  const currentPrice = finiteNumber(stock?.currentPrice);
  if (!Number.isFinite(perShare) || perShare <= 0 || !Number.isFinite(currentPrice) || currentPrice <= 0) {
    return null;
  }

  return {
    cashDividendTotalCny: perShare * fx(dividendCurrency(stock)),
    marketCapCny: currentPrice * fx(stock?.currency || dividendCurrency(stock))
  };
}

function calculatedDividendYield(stock) {
  const inputs = dividendYieldInputs(stock);
  if (!inputs || !Number.isFinite(inputs.marketCapCny) || inputs.marketCapCny <= 0) return null;
  return inputs.cashDividendTotalCny / inputs.marketCapCny;
}

function calculatedShareholderReturnYield(stock) {
  const inputs = dividendYieldInputs(stock);
  if (!inputs || !Number.isFinite(inputs.marketCapCny) || inputs.marketCapCny <= 0) return null;
  const buybackAmount = finiteNumber(stock?.dividend?.buybackAmount);
  const buybackCny = Number.isFinite(buybackAmount) && buybackAmount > 0
    ? buybackAmount * fx(buybackCurrency(stock))
    : 0;
  return (inputs.cashDividendTotalCny + buybackCny) / inputs.marketCapCny;
}

function dividendAnnualCashLocal(stock) {
  const dividend = stock?.dividend;
  if (!dividend) return null;
  const perShare = finiteNumber(dividend.dividendPerShare);
  if (Number.isFinite(perShare) && Number.isFinite(stock.shares) && stock.shares > 0) {
    return perShare * stock.shares;
  }
  return finiteNumber(dividend.estimatedAnnualCash);
}

function dividendAnnualCashCny(stock) {
  const localCash = dividendAnnualCashLocal(stock);
  if (!Number.isFinite(localCash)) return 0;
  return localCash * fx(dividendCurrency(stock));
}

function dividendSummary(positions) {
  const items = positions
    .map((position) => ({
      position,
      annualCashCny: dividendAnnualCashCny(position),
      reliability: dividendReliability(position)
    }))
    .filter((item) => item.annualCashCny > 0)
    .sort((a, b) => b.annualCashCny - a.annualCashCny);
  const annualCashCny = items.reduce((sum, item) => sum + item.annualCashCny, 0);
  const highRiskCashCny = items
    .filter((item) => item.reliability.value === "risk")
    .reduce((sum, item) => sum + item.annualCashCny, 0);
  return {
    annualCashCny,
    highRiskCashCny,
    topContributor: items.find((item) => item.reliability.value !== "risk") ?? items[0] ?? null
  };
}

function displayDividendRatio(value) {
  return Number.isFinite(value) ? percent(value * 100, false) : "-";
}

function financialAnnuals(stock) {
  return Array.isArray(stock?.financials?.annual) ? stock.financials.annual : [];
}

function latestAnnualFinancial(stock) {
  return financialAnnuals(stock)[0] ?? {};
}

function financialValuation(stock) {
  return stock?.financials?.valuation ?? {};
}

function financialRatio(value, fallback = "-") {
  const number = finiteNumber(value);
  return Number.isFinite(number) ? percent(number * 100, false) : fallback;
}

function financialMultiple(value, fallback = "-") {
  const number = finiteNumber(value);
  return Number.isFinite(number) ? `${number.toFixed(number >= 10 ? 1 : 2)}x` : fallback;
}

function recentAverage(stock, key, count = 5) {
  const values = financialAnnuals(stock)
    .slice(0, count)
    .map((item) => finiteNumber(item?.[key]))
    .filter(Number.isFinite);
  if (!values.length) return NaN;
  return values.reduce((sum, value) => sum + value, 0) / values.length;
}

function positiveRecordRatio(stock, key, count = 5) {
  const values = financialAnnuals(stock)
    .slice(0, count)
    .map((item) => finiteNumber(item?.[key]))
    .filter(Number.isFinite);
  if (!values.length) return NaN;
  return values.filter((value) => value > 0).length / values.length;
}

function compoundGrowth(stock, key, count = 5) {
  const values = financialAnnuals(stock)
    .slice(0, count)
    .map((item) => finiteNumber(item?.[key]))
    .filter((value) => Number.isFinite(value) && value > 0);
  if (values.length < 2) return NaN;
  const current = values[0];
  const oldest = values[values.length - 1];
  if (oldest <= 0) return NaN;
  return Math.pow(current / oldest, 1 / (values.length - 1)) - 1;
}

function financialAmount(value, currencyCode = "") {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return "-";
  const sign = number < 0 ? "-" : "";
  const abs = Math.abs(number);
  const code = currencyCode ? `${currencyCode} ` : "";
  if (abs >= 1e12) return `${sign}${code}${(abs / 1e12).toFixed(2)}万亿`;
  if (abs >= 1e8) return `${sign}${code}${(abs / 1e8).toFixed(2)}亿`;
  if (abs >= 1e4) return `${sign}${code}${(abs / 1e4).toFixed(2)}万`;
  return `${sign}${code}${abs.toFixed(2)}`;
}

function rangeText(range, formatter = financialMultiple) {
  if (!range) return "-";
  return `${formatter(range.min)} / ${formatter(range.median)} / ${formatter(range.max)}`;
}

function dividendReliability(stock) {
  const explicit = String(stock?.dividend?.reliability ?? "").trim().toLowerCase();
  if (["stable", "review", "risk"].includes(explicit)) {
    if (explicit === "stable") return { value: "stable", text: "稳定", tone: "strong" };
    if (explicit === "risk") return { value: "risk", text: "高风险", tone: "risk" };
    return { value: "review", text: "需复核", tone: "watch" };
  }

  const dividend = stock?.dividend;
  if (!dividend) return { value: "review", text: "需复核", tone: "watch" };
  const fiscalYear = String(dividend.fiscalYear ?? "").trim();
  const cashDividendTotal = finiteNumber(dividend.cashDividendTotal);
  const marketCap = finiteNumber(stock?.marketCap);
  if (hasMajorRisk(stock) || valuationConfidence(stock) === "low") {
    return { value: "risk", text: "高风险", tone: "risk" };
  }
  if (!Number.isFinite(cashDividendTotal) || cashDividendTotal <= 0 || !Number.isFinite(marketCap) || marketCap <= 0 || /^TTM/i.test(fiscalYear)) {
    return { value: "review", text: "需复核", tone: "watch" };
  }
  return { value: "stable", text: "稳定", tone: "strong" };
}

function marketKind(stock) {
  const symbol = normalizeSymbol(stock?.symbol);
  if (symbol.endsWith(".HK")) return "HK";
  if (symbol.endsWith(".SH") || symbol.endsWith(".SZ") || symbol.endsWith(".SS")) return "A";
  return String(stock?.currency ?? "").toUpperCase() === "HKD" ? "HK" : "A";
}

function shareholderReturnTarget(stock) {
  return marketKind(stock) === "HK" ? HK_SHARE_SHAREHOLDER_RETURN_TARGET : A_SHARE_SHAREHOLDER_RETURN_TARGET;
}

function forecastDividendYield(stock) {
  const dividend = stock?.dividend;
  const explicit = finiteNumber(dividend?.forecastYield);
  if (Number.isFinite(explicit) && explicit > 0) return explicit;
  const perShare = finiteNumber(dividend?.forecastPerShare);
  const currentPrice = finiteNumber(stock?.currentPrice);
  if (!Number.isFinite(perShare) || perShare <= 0 || !Number.isFinite(currentPrice) || currentPrice <= 0) return null;
  const forecastCurrency = String(dividend?.forecastCurrency || dividendCurrency(stock)).toUpperCase();
  return (perShare * fx(forecastCurrency)) / (currentPrice * fx(stock?.currency || forecastCurrency));
}

function dividendShield(stock) {
  const trailing = calculatedDividendYield(stock);
  const shareholderReturn = calculatedShareholderReturnYield(stock);
  const forecast = forecastDividendYield(stock);
  const value = Number.isFinite(shareholderReturn) ? shareholderReturn : null;
  const source = Number.isFinite(shareholderReturn) ? "最近财年综合回报" : "未记录综合回报";
  const target = shareholderReturnTarget(stock);
  return {
    trailing,
    shareholderReturn,
    forecast,
    value,
    source,
    target,
    passed: Number.isFinite(value) && value >= target
  };
}

function dcfMargin(stock) {
  return calculatedMarginOfSafety(stock) ?? finiteNumber(stock?.marginOfSafety);
}

function amountCurrency(stock, fallback = "") {
  return String(fallback || stock?.currency || "CNY").toUpperCase();
}

function toCnyAmount(value, currencyCode) {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return null;
  return number * fx(String(currencyCode || "CNY").toUpperCase());
}

function companyMarketCapCny(stock) {
  const marketCap = finiteNumber(stock?.marketCap);
  if (!Number.isFinite(marketCap) || marketCap <= 0) return null;
  return toCnyAmount(marketCap, marketCapCurrency(stock));
}

function latestFinancialCurrency(stock) {
  return amountCurrency(stock, latestAnnualFinancial(stock)?.currency);
}

function netCashHaircut(stock) {
  if (hasMajorRisk(stock)) return 0;
  const explicit = finiteNumber(stock?.netCash?.haircut);
  if (Number.isFinite(explicit)) return clamp(explicit, 0, 1);
  const reliability = dividendReliability(stock).value;
  if (reliability === "stable") return 1;
  const text = riskText(stock);
  if (reliability === "risk" || /周期|下滑|亏损|补助|地产链|现金流为负|审计|调查/.test(text)) return 0.4;
  return 0.7;
}

function netCashProfile(stock) {
  const profile = stock?.netCash ?? {};
  const profileCurrency = amountCurrency(stock, profile.currency);
  const rawNetCash = finiteNumber(profile.netCash);
  const cash = finiteNumber(profile.cashAndShortInvestments);
  const debt = finiteNumber(profile.interestBearingDebt);
  const netCash = Number.isFinite(rawNetCash)
    ? rawNetCash
    : Number.isFinite(cash) && Number.isFinite(debt)
      ? cash - debt
      : null;
  const haircut = netCashHaircut(stock);
  const adjustedLocal = finiteNumber(profile.adjustedNetCash) ?? (Number.isFinite(netCash) ? netCash * haircut : null);
  const adjustedCny = toCnyAmount(adjustedLocal, profileCurrency);
  const marketCapCny = companyMarketCapCny(stock);
  const latest = latestAnnualFinancial(stock);
  const financialCurrency = latestFinancialCurrency(stock);
  const netProfitCny = toCnyAmount(latest.netProfit, financialCurrency);
  const shareholderFcf = finiteNumber(profile.shareholderFcf);
  const shareholderFcfCurrency = String(profile.shareholderFcfCurrency || profile.currency || financialCurrency).toUpperCase();
  const fcfCurrency = Number.isFinite(shareholderFcf) ? shareholderFcfCurrency : financialCurrency;
  const fcfLocal = Number.isFinite(shareholderFcf) ? shareholderFcf : finiteNumber(latest.freeCashFlow);
  const fcfCny = toCnyAmount(fcfLocal, fcfCurrency);
  const exCashMarketCapCny = Number.isFinite(marketCapCny) && Number.isFinite(adjustedCny)
    ? Math.max(marketCapCny - Math.max(adjustedCny, 0), 0)
    : null;
  const computedExCashPe = Number.isFinite(exCashMarketCapCny) && Number.isFinite(netProfitCny) && netProfitCny > 0
    ? exCashMarketCapCny / netProfitCny
    : null;
  const computedExCashPfcf = Number.isFinite(exCashMarketCapCny) && Number.isFinite(fcfCny) && fcfCny > 0
    ? exCashMarketCapCny / fcfCny
    : null;
  const computedPfcf = Number.isFinite(marketCapCny) && marketCapCny > 0 && Number.isFinite(fcfCny) && fcfCny > 0
    ? marketCapCny / fcfCny
    : null;
  const exCashPfcf = finiteNumber(profile.exCashPfcf) ?? computedExCashPfcf;
  const pfcf = finiteNumber(profile.pfcf) ?? computedPfcf;
  const computedFcfYield = Number.isFinite(marketCapCny) && marketCapCny > 0 && Number.isFinite(fcfCny)
    ? fcfCny / marketCapCny
    : null;
  return {
    cash,
    debt,
    netCash,
    currency: profileCurrency,
    haircut,
    adjustedLocal,
    adjustedCny,
    marketCapCny,
    exCashPe: finiteNumber(profile.exCashPe) ?? computedExCashPe,
    exCashPfcf,
    pfcf,
    fcfMultiple: exCashPfcf ?? pfcf,
    fcfYield: finiteNumber(profile.fcfYield) ?? computedFcfYield,
    fcfLocal,
    fcfCurrency,
    fcfBasis: String(profile.shareholderFcfBasis || (Number.isFinite(shareholderFcf) ? "普通股东 FCF" : "合并 FCF")).trim(),
    shareholderFcf,
    consolidatedFcf: finiteNumber(profile.consolidatedFcf),
    minorityFcfAdjustment: finiteNumber(profile.minorityFcfAdjustment),
    fcfRecord: positiveRecordRatio(stock, "freeCashFlow"),
    fcfPositiveYears: finiteNumber(profile.fcfPositiveYears),
    reason: String(profile.haircutReason || profile.note || "").trim()
  };
}

function ownerAuditStatusMeta(status) {
  return OWNER_AUDIT_STATUS[status] ?? OWNER_AUDIT_STATUS.review;
}

function ownerAuditScoreMeta(score, hasAudit) {
  if (!hasAudit) return { status: "review", text: "长期股东待评分", tone: "watch", grade: "待补充" };
  if (score >= 85) return { status: "pass", text: "长期股东强", tone: "strong", grade: "优秀" };
  if (score >= OWNER_AUDIT_SCORE_TARGET) return { status: "pass", text: "长期股东达标", tone: "strong", grade: "达标" };
  if (score >= 60) return { status: "review", text: "长期股东观察", tone: "watch", grade: "观察" };
  return { status: "fail", text: "长期股东偏弱", tone: "risk", grade: "偏弱" };
}

function normalizeOwnerAuditStatus(status) {
  const value = String(status ?? "").trim().toLowerCase();
  return ["pass", "review", "fail"].includes(value) ? value : "review";
}

function isCorruptedAuditNote(note) {
  const text = String(note ?? "").trim();
  return /\?{3,}/.test(text) && !/[\u4e00-\u9fa5]/.test(text);
}

function ownerAuditProfile(stock) {
  const raw = stock?.ownerCashFlowAudit;
  const hasAudit = raw && OWNER_AUDIT_FIELDS.some(({ key }) => {
    const item = raw[key] ?? {};
    return String(item.status ?? "").trim() || String(item.note ?? "").trim();
  });
  const items = OWNER_AUDIT_FIELDS.map((field) => {
    const source = raw?.[field.key] ?? {};
    const status = hasAudit ? normalizeOwnerAuditStatus(source.status) : "review";
    const rawNote = String(source.note ?? "").trim();
    const noteCorrupted = isCorruptedAuditNote(rawNote);
    const note = noteCorrupted ? "" : rawNote;
    const weight = finiteNumber(field.weight) ?? 0;
    const statusScore = hasAudit ? OWNER_AUDIT_STATUS_SCORE[status] ?? OWNER_AUDIT_STATUS_SCORE.review : 0;
    const score = Math.round(weight * statusScore);
    return { ...field, status, note, rawNote, noteCorrupted, score, maxScore: weight, ...ownerAuditStatusMeta(status) };
  });
  const maxScore = items.reduce((sum, item) => sum + item.maxScore, 0);
  const rawScore = items.reduce((sum, item) => sum + item.score, 0);
  const score = maxScore > 0 ? Math.round((rawScore / maxScore) * 100) : 0;
  const failItems = items.filter((item) => item.status === "fail");
  const reviewItems = items.filter((item) => item.status === "review");
  const valuationSystem = items.find((item) => item.key === "valuationSystemRisk");
  const corruptedNotes = items.filter((item) => item.noteCorrupted).length;
  const meta = ownerAuditScoreMeta(score, hasAudit);
  const blockers = [];
  if (!hasAudit) {
    blockers.push("长期股东评分待补充");
  } else if (score < OWNER_AUDIT_SCORE_TARGET) {
    blockers.push(`长期股东评分 ${score}/100，低于${OWNER_AUDIT_SCORE_TARGET}`);
  }
  if (valuationSystem?.status === "fail") {
    blockers.push("估值体系失败");
  }
  if (score < OWNER_AUDIT_SCORE_TARGET) {
    blockers.push(...failItems.slice(0, 2).map((item) => `${item.label}0分`));
    if (!failItems.length) {
      blockers.push(...reviewItems.filter((item) => item.core).slice(0, 2).map((item) => `${item.label}待复核`));
    }
  }
  return {
    status: meta.status,
    text: meta.text,
    tone: meta.tone,
    grade: meta.grade,
    score,
    rawScore,
    maxScore,
    hasAudit,
    corruptedNotes,
    items,
    blockers,
    valuationSystemFailed: valuationSystem?.status === "fail"
  };
}

function strategyProfile(stock) {
  const margin = dcfMargin(stock);
  const shield = dividendShield(stock);
  const confidence = valuationConfidence(stock);
  const majorRisk = hasMajorRisk(stock);
  const dividendRisk = dividendReliability(stock).value === "risk";
  const dcfPassed = Number.isFinite(margin) && margin >= MAIN_DCF_MARGIN_TARGET;
  const ownerAudit = ownerAuditProfile(stock);
  const structuralRisk = ownerAudit.valuationSystemFailed;
  const mainPassed = shield.passed && dcfPassed && ownerAudit.status === "pass" && !majorRisk && !structuralRisk && confidence !== "low" && !dividendRisk;
  const netCash = netCashProfile(stock);
  const peLimit = marketKind(stock) === "HK" ? HK_SHARE_EX_CASH_PE_MAX : A_SHARE_EX_CASH_PE_MAX;
  const fcfOk = (Number.isFinite(netCash.exCashPfcf) && netCash.exCashPfcf <= peLimit) ||
    (Number.isFinite(netCash.fcfYield) && netCash.fcfYield >= 0.1);
  const cigarPassed = !majorRisk && !structuralRisk &&
    Number.isFinite(netCash.adjustedCny) && netCash.adjustedCny > 0 &&
    Number.isFinite(netCash.exCashPe) && netCash.exCashPe <= peLimit &&
    fcfOk;

  let bucket = "transition";
  let status = "过渡观察";
  let tone = "watch";
  if (majorRisk || structuralRisk || confidence === "low" || dividendRisk) {
    bucket = "excluded";
    status = "风险排除";
    tone = "risk";
  } else if (mainPassed) {
    bucket = "main";
    status = "主策略达标";
    tone = "buy";
  } else if (cigarPassed) {
    bucket = "cigar";
    status = "辅策略烟蒂";
    tone = "safe";
  } else if (shield.passed || dcfPassed || Number.isFinite(netCash.exCashPe)) {
    status = "过渡观察";
    tone = "wait";
  }

  const blockers = [];
  if (ownerAudit.status !== "pass") blockers.push(...ownerAudit.blockers.slice(0, 2));
  if (!shield.passed) blockers.push(`综合回报未达 ${percent(shield.target * 100, false)}`);
  if (!dcfPassed) blockers.push(`安全边际未达 ${percent(MAIN_DCF_MARGIN_TARGET * 100, false)}`);
  if (majorRisk) blockers.push("重大风险未解除");
  if (structuralRisk) blockers.push("估值体系改变风险");
  if (confidence === "low") blockers.push("估值低可信");
  if (!cigarPassed && Number.isFinite(netCash.exCashPe)) blockers.push(`烟蒂PE需≤${peLimit}x且FCF达标`);

  return {
    bucket,
    status,
    tone,
    margin,
    shield,
    confidence,
    ownerAudit,
    netCash,
    peLimit,
    mainPassed,
    cigarPassed,
    blockers
  };
}

function strategyBucketLabel(bucket) {
  if (bucket === "main") return "主策略";
  if (bucket === "cigar") return "辅策略";
  if (bucket === "excluded") return "风险排除";
  return "过渡观察";
}

function strategyUniverseItems(positions) {
  return decisionUniverse(positions)
    .map((stock) => ({ stock, strategy: strategyProfile(stock) }))
    .sort((a, b) => {
      const order = { main: 0, cigar: 1, transition: 2, excluded: 3 };
      return (order[a.strategy.bucket] ?? 9) - (order[b.strategy.bucket] ?? 9) ||
        (b.strategy.shield.value ?? -Infinity) - (a.strategy.shield.value ?? -Infinity) ||
        (b.strategy.margin ?? -Infinity) - (a.strategy.margin ?? -Infinity) ||
        a.stock.name.localeCompare(b.stock.name, "zh-CN");
    });
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

function renderQuoteUpdateStatus(positions) {
  const stocks = auditUniverse(positions);
  const referenceDate = quoteReferenceDate(stocks);
  if (!referenceDate) {
    setQuoteUpdateStatus("行情待更新");
    return;
  }

  const missingQuotes = stocks.filter((stock) => !stock.currentPriceDate || !Number.isFinite(finiteNumber(stock.currentPrice)) || finiteNumber(stock.currentPrice) <= 0);
  const staleQuotes = stocks.filter((stock) => {
    const lagDays = dateDiffDays(stock.currentPriceDate, referenceDate);
    return Number.isFinite(lagDays) && lagDays > 2;
  });
  const missingPreviousClose = stocks.filter((stock) => !stock.previousCloseDate || !Number.isFinite(finiteNumber(stock.previousClose)) || finiteNumber(stock.previousClose) <= 0);

  if (missingQuotes.length || staleQuotes.length || missingPreviousClose.length) {
    const issues = [];
    if (missingQuotes.length) issues.push(`${missingQuotes.length} 个缺行情`);
    if (staleQuotes.length) issues.push(`${staleQuotes.length} 个行情落后`);
    if (missingPreviousClose.length) issues.push(`${missingPreviousClose.length} 个缺昨收`);
    setQuoteUpdateStatus(`行情已至 ${referenceDate}，${issues.join("，")}`, "error");
    return;
  }

  setQuoteUpdateStatus(`行情已更新至 ${referenceDate}`, "success");
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

function findRawStock(symbol) {
  const normalized = normalizeSymbol(symbol);
  return (
    state.holdings.find((stock) => normalizeSymbol(stock.symbol) === normalized) ||
    state.candidates.find((stock) => normalizeSymbol(stock.symbol) === normalized) ||
    null
  );
}

function previewValue(value, formatter = (item) => displayText(item)) {
  if (value === null || value === undefined || value === "") return "空";
  return formatter(value);
}

function researchPreviewDiff(existing, research) {
  if (!existing) return [];
  const isEventUpdate = research.updateType === "eventUpdate";
  const updates = isEventUpdate ? (research.updates ?? {}) : research;
  const valuation = updates.valuation ?? {};
  const quality = updates.quality ?? {};
  const nextLevels = priceLevels({
    targetBuyPrice: valuation.targetBuyPrice,
    intrinsicValue: valuation.intrinsicValue
  });
  const currentLevels = priceLevels(existing);
  const rows = [];
  const push = (label, before, after, formatter) => {
    const beforeText = previewValue(before, formatter);
    const afterText = previewValue(after, formatter);
    if (beforeText !== afterText) rows.push({ label, beforeText, afterText });
  };

  if (!isEventUpdate || valuation.intrinsicValue !== undefined) {
    push("内在价值", existing.intrinsicValue, valuation.intrinsicValue, (value) => currency(Number(value), research.currency || existing.currency || "CNY"));
  }
  if (!isEventUpdate || valuation.fairValueRange !== undefined) push("公允区间", existing.fairValueRange, valuation.fairValueRange);
  if (!isEventUpdate || valuation.intrinsicValue !== undefined || valuation.targetBuyPrice !== undefined) {
    push("首买价", currentLevels.initialBuyPrice, nextLevels.initialBuyPrice, (value) => currency(Number(value), research.currency || existing.currency || "CNY"));
  }
  if (!isEventUpdate || quality.totalScore !== undefined) push("质量总分", existing.qualityScore, quality.totalScore, (value) => String(value));
  if (!isEventUpdate || updates.status !== undefined) push("达标状态", existing.status, updates.status);
  if (!isEventUpdate || updates.action !== undefined) push("最终动作", existing.action, updates.action);
  if (!isEventUpdate || updates.risk !== undefined) push("主要风险", existing.risk, updates.risk);
  if (isEventUpdate && updates.notesAppend) rows.push({ label: "研究脉络", beforeText: "保留原 notes", afterText: updates.notesAppend });
  return rows;
}

function renderResearchPreview(result, imported = false) {
  const research = result.research ?? {};
  const isEventUpdate = research.updateType === "eventUpdate";
  const updates = isEventUpdate ? (research.updates ?? {}) : research;
  const valuation = updates.valuation ?? {};
  const levels = priceLevels({
    targetBuyPrice: valuation.targetBuyPrice,
    intrinsicValue: valuation.intrinsicValue
  });
  const dividend = updates.dividend ?? {};
  const quality = updates.quality ?? {};
  const audit = ownerAuditProfile(isEventUpdate ? { ownerCashFlowAudit: updates.ownerCashFlowAudit } : research);
  const warnings = result.warnings ?? [];
  const plan = result.plan ?? [];
  const existing = findRawStock(research.symbol);
  const diffRows = researchPreviewDiff(existing, research);
  const changedFields = result.changedFields ?? [];
  const event = research.event ?? {};
  const impact = research.impact ?? {};
  const targetPlan = plan.find((item) => {
    const itemSymbol = String(item.symbol ?? "").toUpperCase();
    return itemSymbol ? itemSymbol === String(research.symbol ?? "").toUpperCase() : item.name === research.name;
  });

  elements.researchPreview.innerHTML = `
    <div class="research-summary">
      <strong>${escapeHTML(targetTypeText(result.targetType))}</strong>
      <span>${escapeHTML(isEventUpdate ? "事件/财报增量更新" : "完整重估")} · ${escapeHTML(result.summary ?? "")}</span>
      ${result.backupPath ? `<small>备份：${escapeHTML(result.backupPath)}</small>` : ""}
    </div>
    ${isEventUpdate ? `
      <div class="research-event-card">
        <div>
          <span>${escapeHTML(event.type || "event")}</span>
          <strong>${escapeHTML(event.title || "未填写事件标题")}</strong>
          <small>${escapeHTML([event.date, event.source].filter(Boolean).join(" · ") || "事件来源待补充")}</small>
        </div>
        <p>${escapeHTML(event.summary || "暂无事件摘要")}</p>
        <small>影响：${escapeHTML([
          impact.thesisChange ? `thesis ${impact.thesisChange}` : "",
          impact.valuationChange ? `valuation ${impact.valuationChange}` : "",
          impact.riskChange ? `risk ${impact.riskChange}` : "",
          impact.actionChange ? `action ${impact.actionChange}` : ""
        ].filter(Boolean).join(" · ") || "未填写影响判断")}</small>
      </div>
    ` : ""}
    <div class="research-preview-grid">
      <div><span>股票</span><strong>${escapeHTML(research.name ?? "-")}</strong><small>${escapeHTML(research.symbol ?? "-")} · ${escapeHTML(research.asOf ?? "-")}</small></div>
      <div><span>安全边际</span><strong>${Number.isFinite(valuation.marginOfSafety) ? percent(valuation.marginOfSafety * 100, false) : "-"}</strong><small>${escapeHTML(valuation.fairValueRange ?? "-")}</small></div>
      <div><span>质量总分</span><strong>${numberText(quality.totalScore)}</strong><small>${numberText(quality.businessModel)}/${numberText(quality.moat)}/${numberText(quality.governance)}/${numberText(quality.financialQuality)}</small></div>
      <div><span>执行排序</span><strong>${targetPlan ? targetPlan.rank : "-"}</strong><small>${escapeHTML(targetPlan ? targetPlan.priority : "未列入 Plan")}</small></div>
      <div><span>买入分层</span><strong>${Number.isFinite(levels.initialBuyPrice) ? currency(levels.initialBuyPrice, research.currency || "CNY") : "-"}</strong><small>观察 ${Number.isFinite(levels.watchPrice) ? currency(levels.watchPrice, research.currency || "CNY") : "-"} · 重仓 ${Number.isFinite(levels.aggressiveBuyPrice) ? currency(levels.aggressiveBuyPrice, research.currency || "CNY") : "-"}</small></div>
      <div><span>股东回报</span><strong>${Number.isFinite(dividend.dividendPerShare) ? currency(dividend.dividendPerShare, dividend.dividendCurrency || research.currency || "CNY") : "-"}</strong><small>${dividend.fiscalYear ? `${escapeHTML(dividend.fiscalYear)} · ` : ""}综合回报按分红+回购/总市值计算</small></div>
      <div><span>长期评分</span><strong>${badge(audit.text, audit.tone)}</strong><small>${escapeHTML(audit.hasAudit ? `${audit.score}/100 · ${audit.grade}` : "待补评分")}</small></div>
    </div>
    ${isEventUpdate ? `
      <div class="research-warnings neutral">
        <span>更新字段：${escapeHTML(changedFields.length ? changedFields.join(" / ") : "仅追加研究记录")}</span>
      </div>
    ` : ""}
    ${warnings.length ? `
      <div class="research-warnings">
        ${warnings.map((item) => `<span>${escapeHTML(item)}</span>`).join("")}
      </div>
    ` : ""}
    ${diffRows.length ? `
      <div class="research-diff">
        <strong>将更新字段</strong>
        ${diffRows.map((row) => `
          <div>
            <span>${escapeHTML(row.label)}</span>
            <small>${escapeHTML(row.beforeText)} → ${escapeHTML(row.afterText)}</small>
          </div>
        `).join("")}
      </div>
    ` : existing ? `<div class="research-diff compact">未发现核心字段变化</div>` : ""}
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
      const shares = finiteNumber(holding.shares) ?? 0;
      const cost = finiteNumber(holding.cost) ?? 0;
      const currentPrice = finiteNumber(holding.currentPrice);
      const previousClose = finiteNumber(holding.previousClose);
      const hasCurrentPrice = Number.isFinite(currentPrice) && currentPrice > 0;
      const marketValueLocal = hasCurrentPrice ? shares * currentPrice : 0;
      const costValueLocal = shares * cost;
      const marketValueCny = marketValueLocal * fx(holding.currency);
      const costValueCny = costValueLocal * fx(holding.currency);
      const pnlCny = hasCurrentPrice ? marketValueCny - costValueCny : null;
      const closeForDayChange = Number.isFinite(previousClose) && previousClose > 0 ? previousClose : currentPrice;
      const dayChange = hasCurrentPrice && Number.isFinite(closeForDayChange)
        ? shares * (currentPrice - closeForDayChange) * fx(holding.currency)
        : null;

      return {
        ...holding,
        shares,
        cost,
        currentPrice,
        previousClose,
        marginOfSafety: calculatedMarginOfSafety({ ...holding, currentPrice }) ?? holding.marginOfSafety,
        marketValueLocal,
        marketValueCny,
        costValueCny,
        pnlCny,
        pnlRate: costValueCny && Number.isFinite(pnlCny) ? (pnlCny / costValueCny) * 100 : null,
        dayChange
      };
    });
}

function syncCash() {
  if (Number.isFinite(state.cash)) return;

  const invested = state.holdings.reduce((sum, holding) => {
    const shares = finiteNumber(holding.shares) ?? 0;
    const currentPrice = finiteNumber(holding.currentPrice) ?? 0;
    return sum + shares * currentPrice * fx(holding.currency);
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
  const highRisk = hasMajorRisk(position);
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

function marginValue(stock) {
  return calculatedMarginOfSafety(stock) ?? finiteNumber(stock?.marginOfSafety);
}

function qualityValue(stock) {
  return finiteNumber(stock?.qualityScore);
}

function scoreBand(value, low, high, fallback = 0) {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return fallback;
  if (high === low) return number >= high ? 1 : 0;
  return clamp((number - low) / (high - low), 0, 1);
}

function inverseScoreBand(value, good, bad, fallback = 0) {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return fallback;
  if (bad === good) return number <= good ? 1 : 0;
  return clamp((bad - number) / (bad - good), 0, 1);
}

function dividendScore(stock) {
  const reliability = dividendReliability(stock).value;
  if (reliability === "stable") return 1;
  if (reliability === "review") return 0.55;
  return 0;
}

function confidenceValue(stock) {
  const confidence = valuationConfidence(stock);
  if (confidence === "high") return 1;
  if (confidence === "medium") return 0.65;
  return 0.2;
}

function isFinancialBusiness(stock) {
  return /银行|保险|证券|券商|金融/.test([stock?.name, stock?.industry].filter(Boolean).join(" "));
}

function balanceSheetScore(stock, fallback = 0.55) {
  if (isFinancialBusiness(stock)) return fallback;
  return inverseScoreBand(latestAnnualFinancial(stock).debtRatio, 0.35, 0.75, fallback);
}

function qualityComposite(stock) {
  const quality = qualityValue(stock);
  if (Number.isFinite(quality)) return clamp(quality / 100, 0, 1);

  const business = finiteNumber(stock?.businessModel);
  const moat = finiteNumber(stock?.moat);
  const governance = finiteNumber(stock?.governance);
  const financial = finiteNumber(stock?.financialQuality);
  const parts = [
    Number.isFinite(business) ? (business / 30) * 0.3 : null,
    Number.isFinite(moat) ? (moat / 25) * 0.25 : null,
    Number.isFinite(governance) ? (governance / 20) * 0.2 : null,
    Number.isFinite(financial) ? (financial / 25) * 0.25 : null
  ].filter((value) => value !== null);
  if (!parts.length) return 0.6;
  return clamp(parts.reduce((sum, value) => sum + value, 0), 0, 1);
}

function weightedScore(parts) {
  return Math.round(clamp(parts.reduce((sum, part) => sum + part.weight * clamp(part.score, 0, 1), 0), 0, 100));
}

function masterTone(status) {
  if (/复盘|不合格|能力圈外|投机|低可信|降权|故事不足/.test(status)) return "risk";
  if (/合格|核心|足够|认可|清晰|可验证/.test(status)) return "strong";
  return "watch";
}

function marksView(stock, totalValue = 0) {
  const margin = marginValue(stock);
  const confidence = valuationConfidence(stock);
  const weight = totalValue && stock.marketValueCny ? stock.marketValueCny / totalValue : 0;
  const dividendRisk = dividendReliability(stock).value === "risk";
  const latestFinancial = latestAnnualFinancial(stock);
  const valuation = financialValuation(stock);
  const debtRatio = finiteNumber(latestFinancial.debtRatio);
  const freeCashFlow = finiteNumber(latestFinancial.freeCashFlow);
  const fcfRecord = positiveRecordRatio(stock, "freeCashFlow");
  const revenueGrowth = finiteNumber(latestFinancial.revenueYoY);
  const profitGrowth = finiteNumber(latestFinancial.netProfitYoY);
  const peg = finiteNumber(valuation.peg);
  const financialStress = (!isFinancialBusiness(stock) && Number.isFinite(debtRatio) && debtRatio > 0.65) ||
    (Number.isFinite(freeCashFlow) && freeCashFlow < 0);
  const highRisk = hasMajorRisk(stock) || confidence === "low" || financialStress;
  const requiredMargin = (confidence === "high" ? 0.25 : confidence === "medium" ? 0.3 : 0.38) +
    (financialStress ? 0.08 : 0) +
    (dividendRisk ? 0.04 : 0);
  const riskCompensationScore = scoreBand(margin, 0, requiredMargin, 0);
  const resilienceScore = (
    balanceSheetScore(stock, 0.5) * 0.4 +
    (Number.isFinite(freeCashFlow) ? (freeCashFlow > 0 ? 1 : 0) : 0.45) * 0.35 +
    (Number.isFinite(fcfRecord) ? fcfRecord : 0.45) * 0.25
  );
  const uncertaintyScore = (
    confidenceValue(stock) * 0.4 +
    (hasMajorRisk(stock) ? 0 : 1) * 0.45 +
    (dividendRisk ? 0 : 1) * 0.15
  );
  const cycleScore = (
    scoreBand(profitGrowth, -0.2, 0.12, 0.45) * 0.4 +
    scoreBand(revenueGrowth, -0.1, 0.1, 0.45) * 0.25 +
    inverseScoreBand(peg, 1, 3, 0.55) * 0.15 +
    riskCompensationScore * 0.2
  );
  const positionScore = weight > 0 ? inverseScoreBand(weight, 0.08, 0.22, 0.75) : 1;
  let score = weightedScore([
    { weight: 30, score: riskCompensationScore },
    { weight: 25, score: resilienceScore },
    { weight: 20, score: uncertaintyScore },
    { weight: 15, score: cycleScore },
    { weight: 10, score: positionScore }
  ]);
  if (!Number.isFinite(margin)) score = Math.min(score, 60);
  else if (margin < requiredMargin * 0.5) score = Math.min(score, 68);
  else if (margin < requiredMargin) score = Math.min(score, 82);
  if (highRisk) score = Math.min(score, 58);

  let status = "等待补偿";
  if (highRisk || dividendRisk || resilienceScore < 0.35) status = "风险复盘";
  else if (riskCompensationScore >= 0.95 && score >= 78) status = "补偿足够";
  else if (score >= 65 && Number.isFinite(margin) && margin >= 0.15) status = "仅观察";

  const support = [
    Number.isFinite(margin) ? `安全边际 ${percent(margin * 100, false)}` : "安全边际待补充",
    `要求补偿 ${percent(requiredMargin * 100, false)}`,
    `${confidenceMeta(stock).text}`,
    isFinancialBusiness(stock) ? "金融股需另看资本充足/拨备" : Number.isFinite(debtRatio) ? `负债率 ${financialRatio(debtRatio)}` : "",
    Number.isFinite(freeCashFlow) ? `FCF ${financialAmount(freeCashFlow, latestFinancial.currency || stock.currency)}` : "",
    weight ? `仓位 ${percent(weight * 100, false)}` : ""
  ].filter(Boolean);
  const against = [
    riskCompensationScore < 1 ? "风险补偿不足" : "",
    highRisk ? "存在重大风险或低可信" : "",
    financialStress ? "杠杆或自由现金流有压力" : "",
    weight > 0.15 ? "仓位过高，新增风险补偿要求提高" : "",
    dividendRisk ? "股息质量需复盘" : ""
  ].filter(Boolean);

  return {
    key: "marks",
    name: "马克斯",
    title: "风险与周期",
    status,
    score,
    tone: masterTone(status),
    support,
    against,
    action: status === "补偿足够" ? "可进入买入复核" : status === "风险复盘" ? "优先复盘或降权" : "等待更高风险补偿"
  };
}

function grahamView(stock) {
  const margin = marginValue(stock);
  const financial = finiteNumber(stock?.financialQuality);
  const dividend = dividendReliability(stock);
  const confidence = valuationConfidence(stock);
  const highRisk = hasMajorRisk(stock) || confidence === "low";
  const latestFinancial = latestAnnualFinancial(stock);
  const debtRatio = finiteNumber(latestFinancial.debtRatio);
  const earningsRecord = positiveRecordRatio(stock, "netProfit");
  const cashFlowRecord = positiveRecordRatio(stock, "operatingCashFlow");
  const defensiveMargin = Number.isFinite(margin) && margin >= (confidence === "high" ? SAFETY_MARGIN_TARGET : 0.33);
  const resilienceScore = (
    (Number.isFinite(financial) ? clamp(financial / 25, 0, 1) : 0.5) * 0.55 +
    balanceSheetScore(stock, 0.5) * 0.45
  );
  const stabilityScore = (
    (Number.isFinite(earningsRecord) ? earningsRecord : 0.4) * 0.55 +
    (Number.isFinite(cashFlowRecord) ? cashFlowRecord : 0.4) * 0.45
  );
  const riskScore = (
    confidenceValue(stock) * 0.45 +
    (hasMajorRisk(stock) ? 0 : 1) * 0.55
  );
  let score = weightedScore([
    { weight: 35, score: scoreBand(margin, 0, 0.3, 0) },
    { weight: 20, score: resilienceScore },
    { weight: 20, score: stabilityScore },
    { weight: 10, score: dividendScore(stock) },
    { weight: 15, score: riskScore }
  ]);
  if (!Number.isFinite(margin)) score = Math.min(score, 62);
  else if (margin < 0.1) score = Math.min(score, 58);
  else if (margin < 0.15) score = Math.min(score, 68);
  else if (margin < SAFETY_MARGIN_TARGET) score = Math.min(score, 79);
  if (highRisk) score = Math.min(score, 55);

  let status = "不合格";
  if (!highRisk && defensiveMargin && score >= 75) status = "防御合格";
  else if (!highRisk && Number.isFinite(margin) && margin >= 0.15 && score >= 60) status = "勉强";

  const support = [
    Number.isFinite(margin) ? `安全边际 ${percent(margin * 100, false)}` : "安全边际待补充",
    Number.isFinite(financial) ? `财务质量 ${financial}/25` : "财务质量待补充",
    Number.isFinite(earningsRecord) ? `近年盈利为正 ${Math.round(earningsRecord * 100)}%` : "",
    isFinancialBusiness(stock) ? "金融股需另看资本充足/拨备" : Number.isFinite(debtRatio) ? `负债率 ${financialRatio(debtRatio)}` : "",
    `股息${dividend.text}`
  ].filter(Boolean);
  const against = [
    defensiveMargin ? "" : "防御型折价不足",
    highRisk ? "财报/治理/低可信风险" : "",
    !isFinancialBusiness(stock) && Number.isFinite(debtRatio) && debtRatio > 0.7 ? "资产负债表防守性不足" : "",
    dividend.value === "risk" ? "分红可靠性不足" : ""
  ].filter(Boolean);

  return {
    key: "graham",
    name: "格雷厄姆",
    title: "便宜与防守",
    status,
    score,
    tone: masterTone(status),
    support,
    against,
    action: status === "防御合格" ? "可列入防御型复核" : status === "勉强" ? "继续等待更低价格" : "不满足防守买入"
  };
}

function buffettView(stock) {
  const quality = qualityValue(stock);
  const moat = finiteNumber(stock?.moat);
  const financial = finiteNumber(stock?.financialQuality);
  const margin = marginValue(stock);
  const confidence = valuationConfidence(stock);
  const highRisk = hasMajorRisk(stock) || confidence === "low";
  const avgRoe = recentAverage(stock, "roe");
  const avgRoic = recentAverage(stock, "roic");
  const fcfRecord = positiveRecordRatio(stock, "freeCashFlow");
  const economicsScore = (
    scoreBand(avgRoe, 0.08, 0.2, 0.45) * 0.5 +
    scoreBand(avgRoic, 0.08, 0.18, 0.45) * 0.5
  );
  const cashScore = (
    (Number.isFinite(fcfRecord) ? fcfRecord : 0.45) * 0.7 +
    scoreBand(latestAnnualFinancial(stock).operatingCashFlowToRevenue, 0.05, 0.25, 0.45) * 0.3
  );
  const durabilityScore = (
    confidenceValue(stock) * 0.4 +
    (hasMajorRisk(stock) ? 0 : 1) * 0.45 +
    balanceSheetScore(stock, 0.55) * 0.15
  );
  let score = weightedScore([
    { weight: 35, score: qualityComposite(stock) },
    { weight: 25, score: economicsScore },
    { weight: 15, score: cashScore },
    { weight: 15, score: durabilityScore },
    { weight: 10, score: scoreBand(margin, 0, 0.2, 0.35) }
  ]);
  if (highRisk) score = Math.min(score, 58);
  if (Number.isFinite(quality) && quality < 75) score = Math.min(score, 65);

  let status = "普通机会";
  if (highRisk || (Number.isFinite(quality) && quality < 75)) status = "能力圈外";
  else if (
    score >= 86 &&
    Number.isFinite(quality) && quality >= 85 &&
    Number.isFinite(avgRoic) && avgRoic >= 0.12 &&
    Number.isFinite(fcfRecord) && fcfRecord >= 0.8 &&
    Number.isFinite(margin) && margin >= 0.2
  ) status = "长期核心";
  else if (score >= 78) status = "好生意等价格";

  const support = [
    Number.isFinite(quality) ? `质量 ${quality}` : "质量待补充",
    Number.isFinite(moat) ? `护城河 ${moat}/25` : "护城河待补充",
    Number.isFinite(avgRoe) ? `近年 ROE ${financialRatio(avgRoe)}` : "",
    Number.isFinite(avgRoic) ? `近年 ROIC ${financialRatio(avgRoic)}` : "",
    Number.isFinite(fcfRecord) ? `FCF 为正 ${Math.round(fcfRecord * 100)}%` : "",
    confidenceMeta(stock).text
  ].filter(Boolean);
  const against = [
    Number.isFinite(margin) && margin < 0.2 ? "价格仍不够舒服" : "",
    highRisk ? "低可信或重大风险" : "",
    Number.isFinite(avgRoe) && avgRoe < 0.1 ? "长期 ROE 不够突出" : "",
    Number.isFinite(fcfRecord) && fcfRecord < 0.6 ? "自由现金流连续性不足" : "",
    Number.isFinite(financial) && financial < 20 ? "现金流/财务质量不够强" : ""
  ].filter(Boolean);

  return {
    key: "buffett",
    name: "巴菲特/芒格",
    title: "好生意与复利",
    status,
    score,
    tone: masterTone(status),
    support,
    against,
    action: status === "长期核心" ? "可等待合理价长期加仓" : status === "好生意等价格" ? "留在核心观察，等价格" : status === "能力圈外" ? "隔离为特殊观察" : "普通候选，不消耗太多注意力"
  };
}

function lynchCategory(stock) {
  const text = [stock?.industry, stock?.action, stock?.status, stock?.risk, stock?.notes].filter(Boolean).join(" ");
  if (hasMajorRisk(stock)) return "困境反转/问题股";
  if (/预期差|反转|复苏|修复|验证|重估|改善/.test(text)) return "困境反转/预期差";
  if (/新能源|AI|机器人|自动化|科技|互联网|智能/.test(text)) return "快速增长/高波动";
  if (/油气|银行|地产|物业|航空|周期|白酒|乳制品|家电/.test(text)) return "稳定增长/周期敏感";
  return "稳定增长/普通成长";
}

function lynchView(stock) {
  const quality = qualityValue(stock);
  const margin = marginValue(stock);
  const confidence = valuationConfidence(stock);
  const highRisk = hasMajorRisk(stock) || confidence === "low";
  const latestFinancial = latestAnnualFinancial(stock);
  const valuation = financialValuation(stock);
  const revenueGrowth = finiteNumber(latestFinancial.revenueYoY);
  const profitGrowth = finiteNumber(latestFinancial.netProfitYoY);
  const revenueCagr = compoundGrowth(stock, "revenue");
  const profitCagr = compoundGrowth(stock, "netProfit");
  const cycleRevenueCagr = compoundGrowth(stock, "revenue", 7);
  const cycleProfitCagr = compoundGrowth(stock, "netProfit", 7);
  const fcfRecord = positiveRecordRatio(stock, "freeCashFlow");
  const debtRatio = finiteNumber(latestFinancial.debtRatio);
  const cashConversion = finiteNumber(latestFinancial.operatingCashFlowToRevenue);
  const peg = finiteNumber(valuation.peg);
  const text = [stock?.industry, stock?.action, stock?.status, stock?.risk, stock?.notes].filter(Boolean).join(" ");
  const hasNegativeGrowthCue = /增长弹性不足|增长乏力|增长放缓|增长承压|缺少增长|低增长|成熟|周期敏感|修复已兑现|低基数|高基数/.test(text);
  const hasNarrativeCue = /复苏|修复|预期差|验证|扩张|海外|AI|机器人|新能源|重估|新业务|份额提升|渗透率|订单/.test(text) ||
    (/增长/.test(text) && !hasNegativeGrowthCue);
  const hasGrowthCue = hasNarrativeCue ||
    (Number.isFinite(revenueGrowth) && revenueGrowth > 0.08) ||
    (Number.isFinite(profitGrowth) && profitGrowth > 0.08) ||
    (Number.isFinite(revenueCagr) && revenueCagr > 0.08) ||
    (Number.isFinite(profitCagr) && profitCagr > 0.08);
  const target = finiteNumber(stock?.targetBuyPrice);
  const price = finiteNumber(stock?.currentPrice);
  const nearTarget = Number.isFinite(target) && target > 0 && Number.isFinite(price) && price <= target * 1.08;
  const category = lynchCategory(stock);
  const growthScore = (
    scoreBand(revenueGrowth, -0.05, 0.15, 0.45) * 0.25 +
    scoreBand(profitGrowth, -0.1, 0.2, 0.45) * 0.25 +
    scoreBand(revenueCagr, 0, 0.12, 0.45) * 0.15 +
    scoreBand(profitCagr, 0, 0.15, 0.45) * 0.15 +
    scoreBand(cycleRevenueCagr, 0, 0.1, 0.45) * 0.1 +
    scoreBand(cycleProfitCagr, 0, 0.12, 0.45) * 0.1
  );
  const storyScore = (
    (hasNarrativeCue ? 1 : hasNegativeGrowthCue ? 0.35 : 0.45) * 0.45 +
    qualityComposite(stock) * 0.3 +
    (hasMajorRisk(stock) ? 0 : 1) * 0.25
  );
  const valuationGrowthScore = (
    inverseScoreBand(peg, 0.8, 2.5, 0.5) * 0.45 +
    scoreBand(margin, 0, 0.2, 0.35) * 0.35 +
    (nearTarget ? 1 : 0.45) * 0.2
  );
  const financialVerificationScore = (
    (Number.isFinite(fcfRecord) ? fcfRecord : 0.45) * 0.35 +
    scoreBand(cashConversion, 0.05, 0.25, 0.45) * 0.3 +
    balanceSheetScore(stock, 0.55) * 0.2 +
    confidenceValue(stock) * 0.15
  );
  const riskExecutionScore = (
    confidenceValue(stock) * 0.35 +
    (hasMajorRisk(stock) ? 0 : 1) * 0.45 +
    (Number.isFinite(profitGrowth) && profitGrowth < 0 ? 0.25 : 1) * 0.2
  );
  let score = weightedScore([
    { weight: 25, score: growthScore },
    { weight: 20, score: storyScore },
    { weight: 25, score: valuationGrowthScore },
    { weight: 15, score: financialVerificationScore },
    { weight: 15, score: riskExecutionScore }
  ]);
  if (!hasGrowthCue) score = Math.min(score, 58);
  if (Number.isFinite(profitGrowth) && profitGrowth < 0 && Number.isFinite(revenueGrowth) && revenueGrowth < 0) score = Math.min(score, 62);
  if (Number.isFinite(revenueGrowth) && revenueGrowth < 0 && Number.isFinite(cycleProfitCagr) && cycleProfitCagr <= 0) score = Math.min(score, 74);
  if (Number.isFinite(peg) && peg > 2.5) score = Math.min(score, 72);
  if (highRisk) score = Math.min(score, 55);

  let status = "继续跟踪";
  if (highRisk) status = "等待验证";
  else if (
    score >= 82 &&
    growthScore >= 0.65 &&
    valuationGrowthScore >= 0.6 &&
    (!Number.isFinite(revenueGrowth) || revenueGrowth >= 0) &&
    (!Number.isFinite(cycleProfitCagr) || cycleProfitCagr > 0)
  ) status = "成长故事清晰";
  else if (score >= 68 && (growthScore >= 0.5 || hasNarrativeCue)) status = "预期差可验证";
  else if (score < 52) status = "故事不足";

  const support = [
    `类型：${category}`,
    hasGrowthCue ? "已有成长/修复线索" : "成长线索待补充",
    Number.isFinite(revenueGrowth) ? `收入增速 ${financialRatio(revenueGrowth)}` : "",
    Number.isFinite(profitGrowth) ? `利润增速 ${financialRatio(profitGrowth)}` : "",
    Number.isFinite(revenueCagr) ? `收入CAGR ${financialRatio(revenueCagr)}` : "",
    Number.isFinite(profitCagr) ? `利润CAGR ${financialRatio(profitCagr)}` : "",
    Number.isFinite(cycleRevenueCagr) ? `长周期收入 ${financialRatio(cycleRevenueCagr)}` : "",
    Number.isFinite(cycleProfitCagr) ? `长周期利润 ${financialRatio(cycleProfitCagr)}` : "",
    Number.isFinite(peg) ? `PEG ${peg.toFixed(2)}` : "",
    Number.isFinite(margin) ? `安全边际 ${percent(margin * 100, false)}` : "安全边际待补充",
    nearTarget ? "价格接近买入纪律" : ""
  ].filter(Boolean);
  const against = [
    highRisk ? "风险或数据可信度仍需验证" : "",
    !hasGrowthCue ? "缺少清晰增长故事" : "",
    hasNegativeGrowthCue ? "文本含低增长或周期敏感线索" : "",
    Number.isFinite(revenueGrowth) && revenueGrowth < 0 ? "收入同比下滑" : "",
    Number.isFinite(cycleProfitCagr) && cycleProfitCagr <= 0 ? "长周期利润未验证持续增长" : "",
    Number.isFinite(profitGrowth) && profitGrowth < 0 ? "利润同比下滑" : "",
    Number.isFinite(peg) && peg > 2.5 ? "成长与估值匹配度不足" : "",
    financialVerificationScore < 0.5 ? "增长缺少现金流或资产负债表验证" : "",
    Number.isFinite(margin) && margin < 0.1 ? "估值赔率不够" : ""
  ].filter(Boolean);

  return {
    key: "lynch",
    name: "彼得林奇",
    title: "成长故事与预期差",
    status,
    score,
    tone: masterTone(status),
    support,
    against,
    action: status === "成长故事清晰" ? "可进入成长假设复核" : status === "预期差可验证" ? "跟踪关键验证点" : status === "等待验证" ? "先等财报或风险落地" : "暂不消耗主要仓位"
  };
}

function masterVotes(stock, totalValue = 0) {
  return {
    graham: grahamView(stock),
    buffett: buffettView(stock),
    lynch: lynchView(stock)
  };
}

function riskCommitteeVote(stock, totalValue = 0) {
  return marksView(stock, totalValue);
}

function masterApproval(vote) {
  if (vote.key === "graham") return vote.status === "防御合格" || vote.status === "勉强";
  if (vote.key === "buffett") return vote.status === "长期核心" || vote.status === "好生意等价格";
  if (vote.key === "lynch") return vote.status === "成长故事清晰" || vote.status === "预期差可验证";
  return false;
}

function consensusCount(stock, totalValue = 0) {
  return Object.values(masterVotes(stock, totalValue)).filter(masterApproval).length;
}

function consensusText(count) {
  if (count >= 3) return "三方共识";
  if (count === 2) return "两方认可";
  if (count === 1) return "单方认可";
  return "无共识";
}

function consensusLabel(stock, totalValue = 0) {
  return consensusText(consensusCount(stock, totalValue));
}

function renderMetrics(positions) {
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const cashValue = finiteNumber(state.cash) ?? 0;
  const totalFunds = totalValue + cashValue;
  const dayChange = positions.reduce((sum, item) => sum + item.dayChange, 0);
  const dividends = dividendSummary(positions);
  const dividendYield = totalValue ? dividends.annualCashCny / totalValue : 0;

  elements.totalFunds.textContent = wholeCurrency(totalFunds);
  elements.totalValue.textContent = wholeCurrency(totalValue);
  if (elements.dayChange) {
    elements.dayChange.textContent = wholeCurrency(dayChange);
    elements.dayChange.className = dayChange >= 0 ? "positive" : "negative";
  }
  if (elements.dayChangeRate) {
    elements.dayChangeRate.textContent = percent(totalValue ? (dayChange / totalValue) * 100 : 0);
    elements.dayChangeRate.className = dayChange >= 0 ? "positive" : "negative";
  }
  elements.annualDividend.textContent = wholeCurrency(dividends.annualCashCny);
  elements.portfolioDividendYield.textContent = dividends.topContributor
    ? `组合股息率 ${percent(dividendYield * 100, false)} · 高风险 ${percent(dividends.annualCashCny ? (dividends.highRiskCashCny / dividends.annualCashCny) * 100 : 0, false)}`
    : "组合股息率 0.00%";
  elements.positionCount.textContent = `${positions.length} 只股票`;
  elements.recordCount.textContent = `${state.holdings.length} 条持仓 · ${state.trades.length} 条交易`;
}

function renderPositions(positions) {
  const filtered = getFilteredPositions(positions);
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);

  if (!filtered.length) {
    elements.positionsBody.innerHTML = `<tr><td colspan="6" class="empty-state">暂无符合条件的持仓</td></tr>`;
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
                <span>${escapeHTML(position.symbol)} · ${position.shares} 股 · 成本 ${currency(position.cost, position.currency)}</span>
              </a>
            </div>
          </td>
          <td data-label="市值/现价">
            <strong>${currency(position.marketValueCny)}</strong>
            <br />
            <small class="quote-date">${currency(position.currentPrice, position.currency)} · ${position.currentPriceDate || "收盘日未知"}</small>
          </td>
          <td data-label="盈亏" class="${pnlClass}">
            ${currency(position.pnlCny)}
            <br />
            <small>${percent(position.pnlRate)}</small>
          </td>
          <td data-label="安全边际">${marginText}</td>
          <td data-label="健康">
            <span class="health-pill ${health.tone}" title="${escapeHTML(health.detail)}">${health.status}</span>
            <br />
            <small class="health-score">${health.score} 分 · 仓位 ${weight.toFixed(1)}% · 质量 ${qualityText}</small>
          </td>
          <td data-label="操作">
            <button class="icon-button edit-holding" data-symbol="${position.symbol}" title="编辑 Excel 信息">✎</button>
          </td>
        </tr>
      `;
    })
    .join("");
}

function renderAllocation(positions) {
  if (!elements.allocationChart) return;
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
      const confidence = stock ? confidenceScore(stock) : 0;
      return { item, stock, margin, quality, confidence };
    })
    .sort((a, b) => {
      const marginA = Number.isFinite(a.margin) ? a.margin : -Infinity;
      const marginB = Number.isFinite(b.margin) ? b.margin : -Infinity;
      return marginB - marginA || b.confidence - a.confidence || (b.quality ?? -Infinity) - (a.quality ?? -Infinity) || a.item.name.localeCompare(b.item.name, "zh-CN");
    });
}

function renderPlanAndCandidates() {
  const plans = sortedPlanItems();
  if (elements.overviewPlanList) {
    elements.overviewPlanList.innerHTML = plans.length
      ? plans.map(({ item, stock, margin }, index) => {
      const confidence = stock ? confidenceMeta(stock) : null;
      return `
      <a class="plan-card" href="${stockHash(stock?.symbol || findSymbolForPlan(item))}">
        <span class="plan-rank">${index + 1}</span>
        <div>
          <strong>${escapeHTML(item.name)}</strong>
          <small>${escapeHTML(item.priority)} · 安全边际 ${Number.isFinite(margin) ? percent(margin * 100, false) : "-"}${confidence ? ` · ${escapeHTML(confidence.text)}` : ""}</small>
        </div>
        <p>${escapeHTML(item.advice)}</p>
      </a>
    `;
      })
      .join("")
      : `<div class="empty-state compact-empty">暂无执行计划</div>`;
  }

  const candidates = sortedCandidates();
  elements.candidateList.innerHTML = candidates.length
    ? candidates.map((item) => {
      const plan = findPlanForStock(item);
      const strategy = strategyProfile(item);
      return `
        <a class="candidate-card" href="${stockHash(item.symbol)}">
          <div class="candidate-head">
            <div>
              <strong>${escapeHTML(item.name)}</strong>
              <span>${escapeHTML(item.symbol)} · ${escapeHTML(item.industry)}</span>
            </div>
            <em>${escapeHTML(strategy.status)}</em>
          </div>
          <div class="candidate-metrics">
            <span>综合回报 <strong>${displayDividendRatio(strategy.shield.value)}</strong></span>
            <span>门槛 <strong>${displayDividendRatio(strategy.shield.target)}</strong></span>
            <span>可信度 <strong>${escapeHTML(confidenceMeta(item).text)}</strong></span>
            <span>最新价 <strong>${Number.isFinite(item.currentPrice) && item.currentPrice > 0 ? currency(item.currentPrice, item.currency) : "-"}</strong></span>
            <span>安全边际 <strong>${Number.isFinite(strategy.margin) ? percent(strategy.margin * 100, false) : "-"}</strong></span>
            <span>长期评分 <strong>${strategy.ownerAudit.hasAudit ? `${strategy.ownerAudit.score}/100` : "-"}</strong></span>
            <span>ex-cash PE <strong>${financialMultiple(strategy.netCash.exCashPe)}</strong></span>
            <span>ex-cash P/FCF <strong>${financialMultiple(strategy.netCash.exCashPfcf)}</strong></span>
          </div>
          <p>${escapeHTML(strategy.blockers.length ? strategy.blockers.join("；") : displayText(plan?.advice, item.action))}</p>
        </a>
      `;
    })
    .join("")
    : `<div class="empty-state compact-empty">当前筛选下暂无候选股</div>`;
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
  const initialBuyPrice = priceLevels(candidate).initialBuyPrice;
  if (!currentPrice || !initialBuyPrice) return null;
  return (currentPrice - initialBuyPrice) / initialBuyPrice;
}

function displayBuyDistance(candidate) {
  const distance = candidateBuyDistance(candidate);
  if (!Number.isFinite(distance)) return "-";
  return distance <= 0 ? "已到买点" : percent(distance * 100, false);
}

function disciplineRankScore(stock) {
  const margin = calculatedMarginOfSafety(stock) ?? finiteNumber(stock.marginOfSafety) ?? -0.5;
  const quality = finiteNumber(stock.qualityScore) ?? 60;
  const distance = candidateBuyDistance(stock);
  const distanceScore = Number.isFinite(distance) ? clamp(1 - Math.max(distance, 0), 0, 1) : 0;
  return quality * 0.45 + clamp(margin, -0.2, 0.5) * 100 * 0.3 + confidenceScore(stock) * 20 + distanceScore * 5;
}

function candidateMatchesFilter(candidate, totalValue) {
  const strategy = strategyProfile(candidate);
  if (candidateFilter === "consensus") return strategy.bucket === "main";
  if (candidateFilter === "margin") return strategy.bucket === "cigar";
  if (candidateFilter === "quality") return strategy.bucket === "transition";
  if (candidateFilter === "nearBuy") {
    const distance = candidateBuyDistance(candidate);
    return Number.isFinite(distance) && distance <= BUY_PROXIMITY;
  }
  if (candidateFilter === "dividend") {
    const yieldValue = strategy.shield.value;
    return Number.isFinite(yieldValue) && yieldValue >= strategy.shield.target;
  }
  if (candidateFilter === "lowConfidence") return strategy.bucket === "excluded";
  return true;
}

function sortedCandidates() {
  const holdingSymbols = new Set(state.holdings.filter((holding) => holding.shares > 0).map((holding) => normalizeSymbol(holding.symbol)));
  const totalValue = computePositions().reduce((sum, item) => sum + item.marketValueCny, 0);
  return state.candidates
    .filter((candidate) => !holdingSymbols.has(normalizeSymbol(candidate.symbol)))
    .filter((candidate) => candidateMatchesFilter(candidate, totalValue))
    .sort((a, b) => {
      const byName = () => a.name.localeCompare(b.name, "zh-CN");
      if (candidateSort === "consensus") {
        const order = { main: 4, cigar: 3, transition: 2, excluded: 1 };
        const strategyA = strategyProfile(a);
        const strategyB = strategyProfile(b);
        return (order[strategyB.bucket] ?? 0) - (order[strategyA.bucket] ?? 0) ||
          (strategyB.shield.value ?? -Infinity) - (strategyA.shield.value ?? -Infinity) ||
          byName();
      }
      if (candidateSort === "discipline") {
        return disciplineRankScore(b) - disciplineRankScore(a) || planRank(a) - planRank(b) || byName();
      }
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
      const strategy = strategyProfile(stock);
      const reasons = [];
      let priority = 4;
      let tone = strategy.tone;

      if (strategy.bucket === "excluded") {
        reasons.push("风险排除");
        priority = 1;
      } else if (strategy.bucket === "main") {
        reasons.push("主策略买点");
        reasons.push(`${strategy.shield.source} ${displayDividendRatio(strategy.shield.value)}`);
        reasons.push(`长期股东 ${strategy.ownerAudit.score}/100`);
        priority = 1;
      } else if (strategy.bucket === "cigar") {
        reasons.push("辅策略烟蒂");
        reasons.push(`ex-cash PE ${financialMultiple(strategy.netCash.exCashPe)}`);
        priority = 2;
      } else {
        reasons.push("过渡观察");
        if (strategy.blockers[0]) reasons.push(strategy.blockers[0]);
        priority = 3;
      }

      return { stock, reasons, priority, tone, marginOfSafety: strategy.margin, strategy };
    })
    .filter((item) => item.reasons.length)
    .sort((a, b) => a.priority - b.priority ||
      (b.strategy.shield.value ?? -Infinity) - (a.strategy.shield.value ?? -Infinity) ||
      (b.strategy.margin ?? -Infinity) - (a.strategy.margin ?? -Infinity) ||
      planRank(a.stock) - planRank(b.stock) ||
      a.stock.name.localeCompare(b.stock.name, "zh-CN"));
}

function renderDecisionItem(signal) {
  const { stock, reasons, tone, marginOfSafety } = signal;
  const sourceText = stock.sourceType === "holding" ? "持仓" : "候选";
  const confidence = confidenceMeta(stock);

  return `
    <a class="decision-item ${tone}" href="${stockHash(stock.symbol)}">
      <div class="decision-title">
        <strong>${escapeHTML(stock.name)}</strong>
        <span>${escapeHTML(sourceText)} · ${escapeHTML(stock.symbol)}</span>
      </div>
      <div class="decision-tags">
        ${reasons.map((reason) => `<span>${escapeHTML(reason)}</span>`).join("")}
        ${badge(confidence.text, confidence.tone)}
      </div>
      <div class="decision-metrics">
        <span>现价 <strong>${localPrice(stock, "currentPrice")}</strong></span>
        <span>综合回报率 <strong>${displayDividendRatio(signal.strategy?.shield.value)}</strong></span>
        <span>回报门槛 <strong>${displayDividendRatio(signal.strategy?.shield.target)}</strong></span>
        <span>安全边际 <strong>${Number.isFinite(marginOfSafety) ? percent(marginOfSafety * 100, false) : "-"}</strong></span>
        <span>长期评分 <strong>${signal.strategy?.ownerAudit.hasAudit ? `${signal.strategy.ownerAudit.score}/100` : "-"}</strong></span>
        <span>ex-cash PE <strong>${financialMultiple(signal.strategy?.netCash.exCashPe)}</strong></span>
        <span>ex-cash P/FCF <strong>${financialMultiple(signal.strategy?.netCash.exCashPfcf)}</strong></span>
      </div>
      <p>${escapeHTML(displayText(stock.action, stock.status))}</p>
    </a>
  `;
}

function renderOpportunityRadar(positions) {
  const signals = buildOpportunitySignals(positions);
  elements.opportunityRadar.innerHTML = signals.length
    ? signals.map(renderDecisionItem).join("")
    : `<div class="empty-state compact-empty">暂无进入买点的标的</div>`;
}

function buildActionConclusion(positions) {
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const totalAssets = totalValue + (finiteNumber(state.cash) ?? 0);
  const cashRatio = totalAssets ? (finiteNumber(state.cash) ?? 0) / totalAssets : 0;
  const signals = buildOpportunitySignals(positions);
  const buySignals = signals.filter(({ strategy }) => strategy.bucket === "main" || strategy.bucket === "cigar");
  const reduceSignals = signals.filter(({ strategy }) => strategy.bucket === "excluded");
  const strategyItems = strategyUniverseItems(positions);
  const mainValue = strategyItems
    .filter((item) => item.stock.sourceType === "holding" && item.strategy.bucket === "main")
    .reduce((sum, item) => sum + (item.stock.marketValueCny ?? 0), 0);
  const cigarValue = strategyItems
    .filter((item) => item.stock.sourceType === "holding" && item.strategy.bucket === "cigar")
    .reduce((sum, item) => sum + (item.stock.marketValueCny ?? 0), 0);
  const transitionValue = strategyItems
    .filter((item) => item.stock.sourceType === "holding" && item.strategy.bucket === "transition")
    .reduce((sum, item) => sum + (item.stock.marketValueCny ?? 0), 0);
  const mainRatio = totalValue ? mainValue / totalValue : 0;
  const cigarRatio = totalValue ? cigarValue / totalValue : 0;
  const transitionRatio = totalValue ? transitionValue / totalValue : 0;
  const dividendRisk = dividendSummary(positions);
  const highRiskDividendRatio = dividendRisk.annualCashCny ? dividendRisk.highRiskCashCny / dividendRisk.annualCashCny : 0;
  const reasons = [];

  reasons.push(`主策略 ${percent(mainRatio * 100, false)} / 目标 ${percent(MAIN_ALLOCATION_TARGET * 100, false)}`);
  reasons.push(`辅策略 ${percent(cigarRatio * 100, false)} / 目标 ${percent(CIGAR_ALLOCATION_TARGET * 100, false)}`);
  if (transitionRatio > 0) reasons.push(`过渡观察 ${percent(transitionRatio * 100, false)}`);
  if (cashRatio >= 0.3) {
    reasons.push(`现金比例 ${percent(cashRatio * 100, false)}，适合等待高赔率`);
  }
  if (highRiskDividendRatio >= 0.15) {
    reasons.push(`高风险股息占 ${percent(highRiskDividendRatio * 100, false)}，需复核现金流质量`);
  }

  if (reduceSignals.length) {
    return {
      tone: "risk",
      status: "优先排除风险",
      detail: `${reduceSignals[0]?.stock.name || "持仓"} 触发重大风险或分红可信问题`,
      reasons: reasons.slice(0, 4)
    };
  }
  if (buySignals.length) {
    const first = buySignals[0];
    return {
      tone: "buy",
      status: "有策略买点",
      detail: `${first.stock.name} 进入${strategyBucketLabel(first.strategy.bucket)}复核`,
      reasons: [first.reasons[0], ...reasons].slice(0, 4)
    };
  }
  if (signals.length) {
    return {
      tone: "watch",
      status: "过渡观察",
      detail: "旧仓不强制卖出，新资金只等综合回报/安全边际或烟蒂条件达标",
      reasons: [signals[0].reasons[0], ...reasons].slice(0, 3)
    };
  }
  return {
    tone: "wait",
    status: "无操作，继续等待",
    detail: "当前没有进入双策略买入区的标的",
    reasons: reasons.slice(0, 3)
  };
}

function renderActionConclusion(positions) {
  const conclusion = buildActionConclusion(positions);
  elements.actionConclusion.className = `executive-hero panel ${conclusion.tone}`;
  elements.actionConclusionStatus.textContent = conclusion.status;
  elements.actionConclusionDetail.textContent = conclusion.detail;
}

function committeeUniverse(positions) {
  return decisionUniverse(positions);
}

function committeeStats(positions) {
  const universe = committeeUniverse(positions);
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const withVotes = universe.map((stock) => ({
    stock,
    votes: masterVotes(stock, totalValue),
    riskVote: riskCommitteeVote(stock, totalValue),
    consensus: consensusCount(stock, totalValue)
  }));
  const holdingsValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const longTermValue = positions
    .filter((stock) => buffettView(stock).status === "长期核心" || buffettView(stock).status === "好生意等价格")
    .reduce((sum, item) => sum + item.marketValueCny, 0);
  const riskReview = withVotes.filter((item) => item.riskVote.status === "风险复盘");
  const defensive = withVotes.filter((item) => item.votes.graham.status === "防御合格");
  const compounders = withVotes.filter((item) => item.votes.buffett.status === "长期核心" || item.votes.buffett.status === "好生意等价格");
  const growthStories = withVotes.filter((item) => item.votes.lynch.status === "成长故事清晰" || item.votes.lynch.status === "预期差可验证");
  return { universe, withVotes, totalValue, holdingsValue, longTermValue, riskReview, defensive, compounders, growthStories };
}

function committeePortfolioStatus(stats) {
  const riskCount = stats.riskReview.filter((item) => item.stock.sourceType === "holding").length;
  const highConsensus = stats.withVotes.filter((item) => item.consensus >= 2).length;
  const buyable = stats.withVotes.some((item) => item.consensus >= 2 && item.riskVote.status === "补偿足够");
  if (riskCount > 0) return "当前组合：风险复盘";
  if (buyable) return "当前组合：可小额加仓";
  if (highConsensus >= 3) return "当前组合：防守等待";
  return "当前组合：继续等待";
}

function executiveActionItems(stats, positions) {
  const items = [];
  const seen = new Set();
  const push = (item) => {
    const symbol = normalizeSymbol(item.symbol);
    if (!symbol || seen.has(symbol) || items.length >= 5) return;
    seen.add(symbol);
    items.push(item);
  };

  stats.riskReview
    .filter((item) => item.stock.sourceType === "holding")
    .sort((a, b) => a.stock.name.localeCompare(b.stock.name, "zh-CN"))
    .forEach((item) => push({
      tone: "risk",
      type: "风险复盘",
      symbol: item.stock.symbol,
      name: item.stock.name,
      meta: `${displayMarginOfSafety(item.stock)} · ${confidenceMeta(item.stock).text}`,
      detail: item.riskVote.action
    }));

  buildOpportunitySignals(positions)
    .filter((signal) => signal.tone === "buy" || signal.tone === "safe")
    .forEach((signal) => push({
      tone: signal.tone === "buy" ? "buy" : "safe",
      type: signal.tone === "buy" ? "买入复核" : "安全边际",
      symbol: signal.stock.symbol,
      name: signal.stock.name,
      meta: `${signal.reasons.slice(0, 2).join("、")} · ${displayMarginOfSafety(signal.stock)}`,
      detail: displayText(signal.stock.action, signal.stock.status)
    }));

  stats.withVotes
    .filter((item) => item.consensus >= 2)
    .sort((a, b) => b.consensus - a.consensus || disciplineRankScore(b.stock) - disciplineRankScore(a.stock))
    .forEach((item) => push({
      tone: "watch",
      type: consensusText(item.consensus),
      symbol: item.stock.symbol,
      name: item.stock.name,
      meta: `${displayMarginOfSafety(item.stock)} · 质量 ${Number.isFinite(item.stock.qualityScore) ? item.stock.qualityScore : "-"}`,
      detail: displayText(item.stock.action, item.stock.status)
    }));

  sortedPlanItems().forEach(({ item, stock, margin }) => {
    const symbol = stock?.symbol || findSymbolForPlan(item);
    push({
      tone: "wait",
      type: "执行计划",
      symbol,
      name: item.name,
      meta: `${item.priority} · 安全边际 ${Number.isFinite(margin) ? percent(margin * 100, false) : "-"}`,
      detail: item.advice
    });
  });

  return items;
}

function renderExecutiveActionItem(item, index) {
  return `
    <a class="next-action-item ${item.tone}" href="${stockHash(item.symbol)}">
      <span class="next-action-rank">${index + 1}</span>
      <div>
        <div class="next-action-head">
          <strong>${escapeHTML(item.name)}</strong>
          <em>${escapeHTML(item.type)}</em>
        </div>
        <small>${escapeHTML(item.symbol)} · ${escapeHTML(item.meta)}</small>
        <p>${escapeHTML(item.detail)}</p>
      </div>
    </a>
  `;
}

function renderCommitteeOverview(positions) {
  const actions = buildOpportunitySignals(positions)
    .slice(0, 6)
    .map((signal) => ({
      tone: signal.tone,
      type: strategyBucketLabel(signal.strategy.bucket),
      symbol: signal.stock.symbol,
      name: signal.stock.name,
        meta: `${displayDividendRatio(signal.strategy.shield.value)}综合回报 · 安全边际 ${Number.isFinite(signal.strategy.margin) ? percent(signal.strategy.margin * 100, false) : "-"} · ${signal.strategy.ownerAudit.text}`,
      detail: signal.strategy.blockers.length ? signal.strategy.blockers.join("；") : displayText(signal.stock.action, signal.stock.status)
    }));
  elements.committeeConsensus.innerHTML = actions.length
    ? actions.map(renderExecutiveActionItem).join("")
    : `<div class="empty-state compact-empty">暂无需要立即处理的动作</div>`;
}

function renderConsensusItem(stock, votes, consensus) {
  return `
    <a class="consensus-item" href="${stockHash(stock.symbol)}">
      <div>
        <strong>${escapeHTML(stock.name)}</strong>
        <span>${escapeHTML(stock.symbol)} · ${escapeHTML(consensusText(consensus))}</span>
      </div>
      <div class="master-tags">
        ${masterTag(votes.graham)}
        ${masterTag(votes.buffett)}
        ${masterTag(votes.lynch)}
      </div>
      <small>${escapeHTML(displayText(stock.action, stock.status))}</small>
    </a>
  `;
}

function masterTag(vote) {
  return `<span class="master-tag ${vote.key} ${vote.tone}" title="${escapeHTML(vote.action)}">${escapeHTML(vote.name)}：${escapeHTML(vote.status)}</span>`;
}

function compactMasterCell(vote) {
  return `<span class="master-mini ${vote.key} ${vote.tone}" title="${escapeHTML(vote.action)}">${escapeHTML(vote.status)}</span>`;
}

function shortReasonText(items, emptyText) {
  const cleanItems = (items ?? []).map((item) => String(item ?? "").trim()).filter(Boolean);
  return cleanItems.length ? cleanItems.slice(0, 3).join("；") : emptyText;
}

function renderMasterVoteCard(vote) {
  return `
    <article class="master-vote-card ${vote.key} ${vote.tone}">
      <div class="master-vote-head">
        <div>
          <p class="eyebrow">${escapeHTML(vote.name)}</p>
          <h3>${escapeHTML(vote.title)}</h3>
        </div>
        <span class="master-score ${vote.tone}" title="本地计算评分，满分100">
          <strong>${vote.score}</strong>
          <small>/100</small>
        </span>
      </div>
      <strong>${escapeHTML(vote.status)}</strong>
      <p>${escapeHTML(vote.action)}</p>
      <dl>
        <div><dt>支持</dt><dd>${escapeHTML(shortReasonText(vote.support, "暂无明确支持项"))}</dd></div>
        <div><dt>反对</dt><dd>${escapeHTML(shortReasonText(vote.against, "暂无主要反对项"))}</dd></div>
      </dl>
    </article>
  `;
}

function renderMasterVotesPanel(stock, totalValue) {
  const strategy = strategyProfile(stock);
  return `
    <section class="panel master-votes-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Strategy Fit</p>
          <h2>策略归属：${escapeHTML(strategy.status)}</h2>
        </div>
      </div>
      <div class="master-vote-grid">
        ${renderStrategyVoteCard("长期股东", strategy.ownerAudit.text, strategy.ownerAudit.tone, [
          strategy.ownerAudit.hasAudit ? `评分 ${strategy.ownerAudit.score}/100，门槛 ${OWNER_AUDIT_SCORE_TARGET}/100` : "评分字段待补充",
          `核心项权重更高：十年需求 / 分红FCF / 估值体系`
        ], strategy.ownerAudit.blockers)}
        ${renderStrategyVoteCard("主策略", strategy.mainPassed ? "达标" : "未达标", strategy.mainPassed ? "buy" : "watch", [
          `综合回报 ${displayDividendRatio(strategy.shield.value)} / 门槛 ${displayDividendRatio(strategy.shield.target)}`,
          `安全边际 ${Number.isFinite(strategy.margin) ? percent(strategy.margin * 100, false) : "-"}`,
          `长期股东评分 ${strategy.ownerAudit.hasAudit ? `${strategy.ownerAudit.score}/100` : "待补"}`,
          `口径 ${strategy.shield.source}`
        ], strategy.blockers)}
        ${renderStrategyVoteCard("辅策略", strategy.cigarPassed ? "烟蒂达标" : "未达标", strategy.cigarPassed ? "safe" : "watch", [
          `调整后净现金 ${financialAmount(strategy.netCash.adjustedCny, "CNY")}`,
          `ex-cash PE ${financialMultiple(strategy.netCash.exCashPe)}`,
          `ex-cash P/FCF ${financialMultiple(strategy.netCash.exCashPfcf)}`
        ], strategy.cigarPassed ? [] : ["净现金、PE或FCF条件不足"])}
        ${renderStrategyVoteCard("风险排除", strategy.bucket === "excluded" ? "排除" : "未触发", strategy.bucket === "excluded" ? "risk" : "strong", [
          confidenceMeta(stock).text,
          dividendReliability(stock).text,
          hasMajorRisk(stock) ? "重大风险词命中" : "未命中重大风险"
        ], strategy.bucket === "excluded" ? strategy.blockers : [])}
        ${renderStrategyVoteCard("组合动作", strategyBucketLabel(strategy.bucket), strategy.tone, [
          strategy.bucket === "main" ? "进入70%主策略池" : strategy.bucket === "cigar" ? "进入30%辅策略池" : strategy.bucket === "transition" ? "旧仓过渡观察" : "等待风险解除",
          displayText(stock.action, stock.status)
        ], [])}
      </div>
    </section>
  `;
}

function renderStrategyVoteCard(name, status, tone, support, against) {
  const blockerText = shortReasonText(against, "无主要阻碍");
  return `
    <article class="master-vote-card ${tone}">
      <div class="master-vote-head">
        <div>
          <p class="eyebrow">${escapeHTML(name)}</p>
          <h3>${escapeHTML(status)}</h3>
        </div>
        ${badge(against.length ? "需处理" : "顺畅", against.length ? "watch" : "strong")}
      </div>
      <div class="strategy-vote-body">
        <div>
          <span>依据</span>
          <p>${escapeHTML(shortReasonText(support, "暂无依据"))}</p>
        </div>
        <div>
          <span>阻碍</span>
          <p>${escapeHTML(blockerText)}</p>
        </div>
      </div>
    </article>
  `;
}

function masterUniverseItems(positions, key) {
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  return committeeUniverse(positions)
    .map((stock) => {
      const votes = masterVotes(stock, totalValue);
      const vote = key === "marks" ? riskCommitteeVote(stock, totalValue) : votes[key];
      const margin = calculatedMarginOfSafety(stock) ?? finiteNumber(stock.marginOfSafety);
      return {
        stock,
        votes,
        vote,
        consensus: Object.values(votes).filter(masterApproval).length,
        margin,
        quality: finiteNumber(stock.qualityScore),
        confidence: confidenceScore(stock)
      };
    })
    .sort((a, b) => {
      if (key === "marks") {
        const riskOrder = { "风险复盘": 0, "等待补偿": 1, "仅观察": 2, "补偿足够": 3 };
        return (riskOrder[a.vote.status] ?? 9) - (riskOrder[b.vote.status] ?? 9) ||
          a.vote.score - b.vote.score ||
          (b.stock.marketValueCny ?? 0) - (a.stock.marketValueCny ?? 0) ||
          a.stock.name.localeCompare(b.stock.name, "zh-CN");
      }
      if (key === "graham") {
        const marginA = Number.isFinite(a.margin) ? a.margin : -Infinity;
        const marginB = Number.isFinite(b.margin) ? b.margin : -Infinity;
        return b.vote.score - a.vote.score || marginB - marginA || b.confidence - a.confidence || a.stock.name.localeCompare(b.stock.name, "zh-CN");
      }
      if (key === "buffett") {
        return b.vote.score - a.vote.score || (b.quality ?? -Infinity) - (a.quality ?? -Infinity) || b.confidence - a.confidence || a.stock.name.localeCompare(b.stock.name, "zh-CN");
      }
      if (key === "lynch") {
        const lynchOrder = { "成长故事清晰": 3, "预期差可验证": 2, "继续跟踪": 1, "等待验证": 0, "故事不足": -1 };
        return (lynchOrder[b.vote.status] ?? 0) - (lynchOrder[a.vote.status] ?? 0) ||
          b.vote.score - a.vote.score ||
          b.confidence - a.confidence ||
          a.stock.name.localeCompare(b.stock.name, "zh-CN");
      }
      return b.vote.score - a.vote.score || b.consensus - a.consensus || a.stock.name.localeCompare(b.stock.name, "zh-CN");
    });
}

function renderMasterSummary(items, key) {
  const approved = key === "marks"
    ? items.filter((item) => item.vote.status === "补偿足够" || item.vote.status === "仅观察").length
    : items.filter((item) => masterApproval(item.vote)).length;
  const holdings = items.filter((item) => item.stock.sourceType === "holding").length;
  const top = items[0];
  const statusCounts = items.reduce((counts, item) => {
    counts.set(item.vote.status, (counts.get(item.vote.status) ?? 0) + 1);
    return counts;
  }, new Map());
  const focusText = [...statusCounts.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 2)
    .map(([status, count]) => `${status} ${count}`)
    .join(" · ");
  const primaryLabel = key === "marks" ? "风险可承受" : key === "graham" ? "防御认可" : key === "buffett" ? "复利认可" : "成长认可";

  return `
    <div class="master-summary-grid">
      <div>
        <span>${primaryLabel}</span>
        <strong>${approved}/${items.length}</strong>
      </div>
      <div>
        <span>覆盖持仓</span>
        <strong>${holdings} 只</strong>
      </div>
      <div>
        <span>当前重点</span>
        <strong>${top ? escapeHTML(top.stock.name) : "-"}</strong>
      </div>
      <div>
        <span>结论分布</span>
        <strong>${escapeHTML(focusText || "待补充")}</strong>
      </div>
    </div>
  `;
}

function renderMasterStockCard(item) {
  const { stock, vote, margin, quality, consensus } = item;
  const sourceText = stock.sourceType === "holding" ? "持仓" : "候选";
  return `
    <a class="master-stock-card ${vote.key} ${vote.tone}" href="${stockHash(stock.symbol)}">
      <div class="master-stock-head">
        <div>
          <strong>${escapeHTML(stock.name)}</strong>
          <span>${escapeHTML(stock.symbol)} · ${escapeHTML(sourceText)} · ${escapeHTML(consensusText(consensus))}</span>
        </div>
        <em>${escapeHTML(vote.status)}</em>
      </div>
      <div class="master-stock-metrics">
        <span>评分 <strong>${vote.score}</strong></span>
        <span>安全边际 <strong>${Number.isFinite(margin) ? percent(margin * 100, false) : "-"}</strong></span>
        <span>质量 <strong>${Number.isFinite(quality) ? quality : "-"}</strong></span>
        <span>共识 <strong>${consensus}/3</strong></span>
      </div>
      <p>${escapeHTML(vote.action)}</p>
      <small>${escapeHTML(shortReasonText(vote.support, "支持项待补充"))}</small>
    </a>
  `;
}

function renderMasterStockList(items) {
  if (!items.length) {
    return `<div class="empty-state compact-empty">暂无可评估标的</div>`;
  }
  return items.map(renderMasterStockCard).join("");
}

function renderMastersPage(positions) {
  if (!elements.grahamSummary || !elements.grahamList || !elements.buffettSummary || !elements.buffettList) return;
  const items = strategyUniverseItems(positions);
  const mainItems = items.filter((item) => item.strategy.bucket === "main");
  const cigarItems = items.filter((item) => item.strategy.bucket === "cigar");
  elements.grahamSummary.innerHTML = renderStrategySummary(mainItems, "主策略达标");
  elements.grahamList.innerHTML = renderStrategyStockList(mainItems, "暂无主策略达标标的");
  elements.buffettSummary.innerHTML = renderStrategySummary(cigarItems, "辅策略烟蒂");
  elements.buffettList.innerHTML = renderStrategyStockList(cigarItems, "暂无辅策略烟蒂标的");
  renderMasterMatrix(positions);
}

function renderStrategySummary(items, label) {
  const holdings = items.filter((item) => item.stock.sourceType === "holding");
  const top = items[0];
  const avgYield = items
    .map((item) => item.strategy.shield.value)
    .filter(Number.isFinite);
  const yieldText = avgYield.length ? percent((avgYield.reduce((sum, value) => sum + value, 0) / avgYield.length) * 100, false) : "-";
  const auditPassed = items.filter((item) => item.strategy.ownerAudit.score >= OWNER_AUDIT_SCORE_TARGET).length;
  return `
    <div class="master-summary-grid">
      <div><span>${escapeHTML(label)}</span><strong>${items.length} 只</strong></div>
      <div><span>其中持仓</span><strong>${holdings.length} 只</strong></div>
      <div><span>评分达标</span><strong>${auditPassed} 只</strong></div>
      <div><span>平均回报</span><strong>${yieldText}</strong></div>
      <div><span>当前重点</span><strong>${top ? escapeHTML(top.stock.name) : "-"}</strong></div>
    </div>
  `;
}

function renderStrategyStockList(items, emptyText) {
  return items.length
    ? items.map(renderStrategyStockCard).join("")
    : `<div class="empty-state compact-empty">${escapeHTML(emptyText)}</div>`;
}

function renderStrategyStockCard(item) {
  const { stock, strategy } = item;
  const sourceText = stock.sourceType === "holding" ? "持仓" : "候选";
  return `
    <a class="master-stock-card ${strategy.tone}" href="${stockHash(stock.symbol)}">
      <div class="master-stock-head">
        <div>
          <strong>${escapeHTML(stock.name)}</strong>
          <span>${escapeHTML(stock.symbol)} · ${escapeHTML(sourceText)} · ${escapeHTML(strategyBucketLabel(strategy.bucket))}</span>
        </div>
        <em>${escapeHTML(strategy.status)}</em>
      </div>
      <div class="master-stock-metrics">
        <span>综合回报 <strong>${displayDividendRatio(strategy.shield.value)}</strong></span>
        <span>门槛 <strong>${displayDividendRatio(strategy.shield.target)}</strong></span>
        <span>安全边际 <strong>${Number.isFinite(strategy.margin) ? percent(strategy.margin * 100, false) : "-"}</strong></span>
        <span>长期评分 <strong>${strategy.ownerAudit.hasAudit ? `${strategy.ownerAudit.score}/100` : "-"}</strong></span>
        <span>ex-cash PE <strong>${financialMultiple(strategy.netCash.exCashPe)}</strong></span>
      </div>
      <p>${escapeHTML(strategy.blockers.length ? strategy.blockers.join("；") : displayText(stock.action, stock.status))}</p>
      <small>${escapeHTML(strategy.netCash.reason || `${strategy.shield.source}口径`)}</small>
    </a>
  `;
}

function renderMasterMatrix(positions) {
  if (!elements.masterMatrix) return;
  const rows = strategyUniverseItems(positions)
    .slice()
    .sort(compareMasterMatrixRows);

  elements.masterMatrix.innerHTML = rows.length
    ? `
      <div class="master-matrix-head">
        <span>标的</span>
        <span>策略归属</span>
        <span>${renderMatrixSortButton("return", "综合回报率")}</span>
        <span>${renderMatrixSortButton("margin", "安全边际")}</span>
        <span>${renderMatrixSortButton("owner", "长期评分")}</span>
        <span>净现金</span>
        <span>FCF估值</span>
      </div>
      ${rows.map(({ stock, strategy }) => `
        <a class="master-matrix-row" href="${stockHash(stock.symbol)}">
          <div>
            <strong>${escapeHTML(stock.name)}</strong>
            <small>${escapeHTML(stock.symbol)} · ${escapeHTML(displayText(stock.industry, "未分类"))}</small>
          </div>
          <span class="master-mini ${strategy.tone}">${escapeHTML(strategy.status)}</span>
          <span class="master-mini ${strategy.shield.passed ? "strong" : "watch"}">${displayDividendRatio(strategy.shield.value)} / ${displayDividendRatio(strategy.shield.target)}</span>
          <span class="master-mini ${Number.isFinite(strategy.margin) && strategy.margin >= MAIN_DCF_MARGIN_TARGET ? "strong" : "watch"}">${dcfMatrixText(stock, strategy.margin)}</span>
          <span class="master-mini ${strategy.ownerAudit.tone}">${strategy.ownerAudit.hasAudit ? `${strategy.ownerAudit.score}/100` : "待评分"}</span>
          <span class="master-mini ${netCashMatrixTone(stock, strategy.netCash)}">${netCashMatrixText(stock, strategy.netCash)}</span>
          <em>${fcfMatrixText(stock, strategy.netCash)}</em>
        </a>
      `).join("")}
    `
    : `<div class="empty-state compact-empty">暂无可对照标的</div>`;
}

function matrixSortValue(item, key) {
  if (key === "return") return finiteNumber(item.strategy?.shield?.value);
  if (key === "owner") return finiteNumber(item.strategy?.ownerAudit?.score);
  return finiteNumber(item.strategy?.margin);
}

function compareMasterMatrixRows(a, b) {
  const valueA = matrixSortValue(a, masterMatrixSort.key);
  const valueB = matrixSortValue(b, masterMatrixSort.key);
  const hasA = Number.isFinite(valueA);
  const hasB = Number.isFinite(valueB);
  if (hasA && hasB && valueA !== valueB) {
    return masterMatrixSort.direction === "asc" ? valueA - valueB : valueB - valueA;
  }
  if (hasA !== hasB) return hasA ? -1 : 1;
  return a.stock.name.localeCompare(b.stock.name, "zh-CN");
}

function renderMatrixSortButton(key, label) {
  const active = masterMatrixSort.key === key;
  const direction = active ? masterMatrixSort.direction : "desc";
  const arrow = active ? (direction === "asc" ? "↑" : "↓") : "↕";
  const nextDirection = active && direction === "desc" ? "asc" : "desc";
  return `
    <button class="matrix-sort-button ${active ? "active" : ""}" type="button" data-master-matrix-sort="${key}" data-next-direction="${nextDirection}" aria-label="${escapeHTML(label)}${direction === "asc" ? "升序" : "降序"}">
      <span>${escapeHTML(label)}</span>
      <strong>${arrow}</strong>
    </button>
  `;
}

function dcfMatrixText(stock, margin) {
  if (!Number.isFinite(margin)) return "-";
  return percent(margin * 100, false);
}

function netCashMatrixTone(stock, netCash) {
  if (!netCashApplicable(stock)) return "watch";
  return Number.isFinite(netCash?.adjustedCny) && netCash.adjustedCny > 0 ? "strong" : "watch";
}

function netCashApplicable(stock) {
  return !/(^|\/)(银行|保险|券商|证券|信托|财富管理)(\/|$)/.test(String(stock?.industry ?? ""));
}

function netCashMatrixText(stock, netCash) {
  if (!netCashApplicable(stock)) return "不适用";
  if (!Number.isFinite(netCash?.adjustedCny)) return "未录入";
  if ((Number.isFinite(netCash?.netCashCny) && netCash.netCashCny <= 0) || netCash.adjustedCny <= 0) return "无净现金";
  return financialAmount(netCash.adjustedCny, "CNY");
}

function fcfMatrixText(stock, netCash) {
  if (!netCashApplicable(stock)) return "不适用";
  return financialMultiple(netCash?.fcfMultiple);
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
  if (!elements.disciplineDashboard) return;
  const investedValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const cashValue = finiteNumber(state.cash) ?? 0;
  const totalAssets = investedValue + cashValue;
  const strategyItems = strategyUniverseItems(positions).filter((item) => item.stock.sourceType === "holding");
  const mainValue = strategyItems.filter((item) => item.strategy.bucket === "main").reduce((sum, item) => sum + (item.stock.marketValueCny ?? 0), 0);
  const cigarValue = strategyItems.filter((item) => item.strategy.bucket === "cigar").reduce((sum, item) => sum + (item.stock.marketValueCny ?? 0), 0);
  const transitionValue = strategyItems.filter((item) => item.strategy.bucket === "transition").reduce((sum, item) => sum + (item.stock.marketValueCny ?? 0), 0);
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
      label: "主策略仓位",
      value: percent(investedValue ? (mainValue / investedValue) * 100 : 0, false),
      detail: `目标 ${percent(MAIN_ALLOCATION_TARGET * 100, false)}`
    },
    {
      label: "辅策略仓位",
      value: percent(investedValue ? (cigarValue / investedValue) * 100 : 0, false),
      detail: `目标 ${percent(CIGAR_ALLOCATION_TARGET * 100, false)}`
    },
    {
      label: "过渡观察",
      value: percent(investedValue ? (transitionValue / investedValue) * 100 : 0, false),
      detail: "旧仓不强制卖出"
    },
    {
      label: "现金比例",
      value: percent(totalAssets ? (cashValue / totalAssets) * 100 : 0, false),
      detail: `现金 ${currency(cashValue)}`
    },
    {
      label: "单股最大仓位",
      value: maxPosition ? percent(totalAssets ? (maxPosition.marketValueCny / totalAssets) * 100 : 0, false) : "-",
      detail: maxPosition ? `${maxPosition.name} · ${currency(maxPosition.marketValueCny)}` : "-"
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

function quoteReferenceDate(stocks) {
  return stocks
    .map((stock) => stock.currentPriceDate)
    .filter(Boolean)
    .sort()
    .at(-1) ?? "";
}

function dateDiffDays(fromDate, toDate) {
  if (!fromDate || !toDate) return null;
  const from = new Date(`${fromDate}T00:00:00`);
  const to = new Date(`${toDate}T00:00:00`);
  if (Number.isNaN(from.getTime()) || Number.isNaN(to.getTime())) return null;
  return Math.round((to.getTime() - from.getTime()) / 86400000);
}

function auditUniverse(positions) {
  const holdingSymbols = new Set(positions.map((position) => normalizeSymbol(position.symbol)));
  const holdings = positions.map((position) => ({ ...position, sourceType: "holding" }));
  const candidates = state.candidates.map((candidate) => ({
    ...candidate,
    sourceType: "candidate",
    duplicateHolding: holdingSymbols.has(normalizeSymbol(candidate.symbol))
  }));
  return [...holdings, ...candidates].filter((stock) => stock.symbol && stock.name);
}

function buildDataQualityIssues(positions) {
  const stocks = auditUniverse(positions);
  const referenceDate = quoteReferenceDate(stocks);
  const issues = [];
  const pushIssue = (tone, stock, title, detail) => {
    issues.push({
      tone,
      symbol: stock?.symbol ?? "",
      name: stock?.name ?? "",
      sourceType: stock?.sourceType ?? "",
      title,
      detail
    });
  };

  stocks.forEach((stock) => {
    const currentPrice = finiteNumber(stock.currentPrice);
    const previousClose = finiteNumber(stock.previousClose);
    const intrinsicValue = finiteNumber(stock.intrinsicValue);
    const qualityScore = finiteNumber(stock.qualityScore);
    const dividend = stock.dividend;
    const dividendYield = calculatedDividendYield(stock);
    const lagDays = dateDiffDays(stock.currentPriceDate, referenceDate);

    if (stock.duplicateHolding) {
      pushIssue("warn", stock, "候选池重复持仓", "该标的已有持仓，候选池展示会被过滤；建议只保留一处维护。");
    }
    if (!Number.isFinite(currentPrice) || currentPrice <= 0) {
      pushIssue("warn", stock, "缺最新价", "会影响市值、安全边际、机会雷达和综合回报率计算。");
    }
    if (!stock.currentPriceDate) {
      pushIssue("warn", stock, "缺收盘日期", "无法判断行情是否为最新收盘价。");
    } else if (Number.isFinite(lagDays) && lagDays > 2) {
      pushIssue("info", stock, "行情日期落后", `当前收盘日 ${stock.currentPriceDate}，组合最新收盘日 ${referenceDate}。`);
    }
    if (!Number.isFinite(previousClose) || previousClose <= 0 || !stock.previousCloseDate) {
      pushIssue("warn", stock, "缺昨收口径", "今日变动会退化或无法计算。");
    }
    if (!Number.isFinite(intrinsicValue) || intrinsicValue <= 0) {
      pushIssue("warn", stock, "缺内在价值", "安全边际和买入价分层无法按公式计算。");
    }
    if (!Number.isFinite(qualityScore)) {
      pushIssue("info", stock, "缺质量分", "持仓健康评分会使用保守默认值。");
    }
    if (!displayText(stock.action, "")) {
      pushIssue("info", stock, "缺执行动作", "详情页和决策日志会缺少可复盘结论。");
    }
    const audit = ownerAuditProfile(stock);
    if (!audit.hasAudit) {
      pushIssue("info", stock, "缺长期股东评分", "主策略评分无法计算；需要补充七项现金流证据。");
    } else if (audit.corruptedNotes > 0) {
      pushIssue("warn", stock, "长期评分说明乱码", `${audit.corruptedNotes} 项说明已损坏，需要从研究材料重建。`);
    }

    if (!dividend) {
      if (stock.sourceType === "holding") {
        pushIssue("info", stock, "缺股息数据", "股息现金流不会计入该持仓。");
      }
      return;
    }

    const fiscalYear = String(dividend.fiscalYear ?? "").trim();
    const perShare = finiteNumber(dividend.dividendPerShare);
    const cashDividendTotal = finiteNumber(dividend.cashDividendTotal);
    const marketCap = finiteNumber(stock.marketCap);
    if (/^TTM/i.test(fiscalYear)) {
      pushIssue("warn", stock, "分红口径仍为 TTM", "应优先改为最新完整财年现金分红总额。");
    }
    if (Number.isFinite(perShare) && perShare > 0 && (!Number.isFinite(cashDividendTotal) || !Number.isFinite(marketCap))) {
      pushIssue("info", stock, "股息率使用每股回退", "缺现金分红总额或总市值时，会用每股分红/现价等价计算。");
    }
    if (Number.isFinite(dividendYield) && dividendYield > 0.1) {
      pushIssue("warn", stock, "股息率异常偏高", "可能混入特别股息或 TTM 派息事件，建议复核财年口径。");
    }
  });

  state.plan.forEach((item) => {
    if (!String(item.symbol ?? "").trim()) {
      pushIssue("info", { name: item.name, sourceType: "plan" }, "Plan 缺股票代码", "可按名称匹配，但代码缺失会降低跳转和更新稳定性。");
    }
  });

  const toneRank = { warn: 0, info: 1, error: -1 };
  return issues.sort((a, b) => {
    return (toneRank[a.tone] ?? 9) - (toneRank[b.tone] ?? 9) || a.name.localeCompare(b.name, "zh-CN") || a.title.localeCompare(b.title, "zh-CN");
  });
}

function renderDataQuality(positions) {
  const issues = buildDataQualityIssues(positions);
  const warningCount = issues.filter((issue) => issue.tone === "warn" || issue.tone === "error").length;
  const infoCount = issues.length - warningCount;
  elements.dataQualityMetric.textContent = warningCount ? `${warningCount} 项` : "通过";
  elements.dataQualityMetric.className = warningCount ? "negative" : "positive";
  elements.dataQualityDetail.textContent = warningCount ? `提醒 ${warningCount} · 备注 ${infoCount}` : infoCount ? `${infoCount} 项备注` : "关键数据完整";

  if (!elements.dataQualityList) return;

  elements.dataQualityList.innerHTML = issues.length
    ? issues.slice(0, 14).map((issue) => {
      const tag = issue.sourceType === "holding" ? "持仓" : issue.sourceType === "candidate" ? "候选" : "Plan";
      const content = `
        <div class="data-quality-head">
          <strong>${escapeHTML(issue.title)}</strong>
          <span>${escapeHTML(tag)}</span>
        </div>
        <div class="data-quality-symbol">${escapeHTML(issue.name || issue.symbol || "未命名")} ${issue.symbol ? `· ${escapeHTML(issue.symbol)}` : ""}</div>
        <small>${escapeHTML(issue.detail)}</small>
      `;
      return issue.symbol
        ? `<a class="data-quality-item ${issue.tone}" href="${stockHash(issue.symbol)}">${content}</a>`
        : `<div class="data-quality-item ${issue.tone}">${content}</div>`;
    }).join("")
    : `<div class="empty-state compact-empty">关键数据完整</div>`;
}

function dividendCashItems(positions) {
  return positions
    .map((position) => {
      const dividend = position.dividend;
      const perShare = finiteNumber(dividend?.dividendPerShare);
      const annualLocal = dividendAnnualCashLocal(position);
      const annualCny = dividendAnnualCashCny(position);
      const dividendYield = calculatedDividendYield(position);
      return {
        position,
        perShare,
        annualLocal,
        annualCny,
        dividendYield,
        reliability: dividendReliability(position),
        currencyCode: dividendCurrency(position),
        fiscalYear: dividend?.fiscalYear
      };
    })
    .filter((item) => Number.isFinite(item.perShare) && item.perShare > 0 && Number.isFinite(item.annualLocal) && item.annualLocal > 0)
    .sort((a, b) => b.annualCny - a.annualCny || a.position.name.localeCompare(b.position.name, "zh-CN"));
}

function renderDividendCashList(positions) {
  const items = dividendCashItems(positions);
  elements.dividendCashList.innerHTML = items.length
    ? items.map((item) => `
      <a class="dividend-cash-item" href="${stockHash(item.position.symbol)}">
        <div class="dividend-cash-head">
          <strong>${escapeHTML(item.position.name)}</strong>
          <span>${escapeHTML(item.position.symbol)} · ${escapeHTML(item.reliability.text)}</span>
        </div>
        <div class="dividend-cash-metrics">
          <span>每股股息 <strong>${currency(item.perShare, item.currencyCode)}</strong></span>
          <span>股息率 <strong>${displayDividendRatio(item.dividendYield)}</strong></span>
          <span>预期年现金 <strong>${currency(item.annualLocal, item.currencyCode)}</strong></span>
        </div>
        <small>${item.annualCny > 0 ? `折人民币 ${currency(item.annualCny)}` : "折人民币 -"}${item.fiscalYear ? ` · ${escapeHTML(item.fiscalYear)}` : ""}${item.reliability.value === "risk" ? " · 高股息不等于低风险" : ""}</small>
      </a>
    `).join("")
    : `<div class="empty-state compact-empty">暂无股息现金流数据，点击更新行情后获取</div>`;
}

function decisionLogTypeText(type) {
  if (type === "research") return "导入分析";
  if (type === "quote") return "更新行情";
  if (type === "trade") return "新增交易";
  if (type === "financials") return "更新财务";
  return "记录";
}

function decisionLogTone(type) {
  if (type === "research") return "research";
  if (type === "quote") return "quote";
  if (type === "trade") return "trade";
  if (type === "financials") return "financials";
  return "";
}

function isMeaningfulDecisionLog(log) {
  if (log.type !== "quote") return true;
  return /触发|进入|离开|跨过|跌破|高于|复盘|减仓|安全边际/.test([log.decision, log.discipline, log.detail].join(" "));
}

function sortedDecisionLogs(symbol = "") {
  const normalizedSymbol = String(symbol ?? "").toUpperCase();
  return [...(state.decisionLogs ?? [])]
    .filter((log) => !normalizedSymbol || String(log.symbol ?? "").toUpperCase() === normalizedSymbol)
    .filter(isMeaningfulDecisionLog)
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
  const logs = sortedDecisionLogs()
    .filter((log) => decisionLogFilter === "all" || log.type === decisionLogFilter)
    .slice(0, 10);
  elements.decisionLogList.innerHTML = renderDecisionLogItems(logs, "暂无符合筛选的决策日志");
}

function renderStockDecisionLogs(stock) {
  return renderDecisionLogItems(sortedDecisionLogs(stock.symbol).slice(0, 8), "暂无该股票的决策日志");
}

function renderResearchUpdatesPanel(stock) {
  const updates = [...(stock.researchUpdates ?? [])].sort((a, b) => String(b.importedAt || b.asOf || "").localeCompare(String(a.importedAt || a.asOf || "")));
  return `
    <section class="panel research-updates-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Research Loop</p>
          <h2>研究更新</h2>
        </div>
      </div>
      <div class="research-update-list">
        ${updates.length ? updates.map((item) => {
          const event = item.event ?? {};
          const impact = item.impact ?? {};
          const impactText = [
            impact.thesisChange ? `thesis ${impact.thesisChange}` : "",
            impact.valuationChange ? `valuation ${impact.valuationChange}` : "",
            impact.riskChange ? `risk ${impact.riskChange}` : "",
            impact.actionChange ? `action ${impact.actionChange}` : ""
          ].filter(Boolean).join(" · ");
          return `
            <article class="research-update-item">
              <div>
                <span>${escapeHTML(event.type || item.updateType || "update")}</span>
                <strong>${escapeHTML(event.title || "未填写事件标题")}</strong>
                <small>${escapeHTML([event.date || item.asOf, event.source].filter(Boolean).join(" · ") || item.importedAt || "-")}</small>
              </div>
              <p>${escapeHTML(event.summary || item.summary || "暂无摘要")}</p>
              <small>${escapeHTML(impactText || "影响判断待补充")}</small>
              <small>更新字段：${escapeHTML((item.changedFields ?? []).join(" / ") || "仅记录")}</small>
              ${item.notesAppend ? `<p class="muted-note">${escapeHTML(item.notesAppend)}</p>` : ""}
            </article>
          `;
        }).join("") : `<div class="empty-state">暂无研究更新记录</div>`}
      </div>
    </section>
  `;
}

function renderDecisionArea(positions) {
  renderActionConclusion(positions);
  renderOpportunityRadar(positions);
  renderDisciplineDashboard(positions);
  renderDataQuality(positions);
  renderDividendCashList(positions);
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

function searchableStocks() {
  const positions = computePositions();
  const seen = new Set();
  const addSource = (items, sourceType) => items
    .filter((stock) => {
      const symbol = normalizeSymbol(stock.symbol);
      if (!symbol || seen.has(symbol)) return false;
      seen.add(symbol);
      return true;
    })
    .map((stock) => ({ ...stock, sourceType }));

  return [
    ...addSource(positions, "持仓"),
    ...addSource(state.candidates, "候选")
  ];
}

function globalSearchMatches(term) {
  const keyword = String(term ?? "").trim().toLowerCase();
  if (!keyword) return [];
  return searchableStocks()
    .filter((stock) => {
      const haystack = [stock.symbol, stock.name, stock.industry].join(" ").toLowerCase();
      return haystack.includes(keyword);
    })
    .sort((a, b) => {
      const exactA = normalizeSymbol(a.symbol).toLowerCase() === keyword || String(a.name).toLowerCase() === keyword;
      const exactB = normalizeSymbol(b.symbol).toLowerCase() === keyword || String(b.name).toLowerCase() === keyword;
      return Number(exactB) - Number(exactA) || a.name.localeCompare(b.name, "zh-CN");
    })
    .slice(0, 8);
}

function renderSearchResults() {
  if (!elements.searchResults) return;
  const matches = globalSearchMatches(elements.searchInput.value);
  elements.searchResults.innerHTML = matches.length
    ? matches.map((stock) => `
      <button type="button" data-search-symbol="${escapeHTML(stock.symbol)}">
        <strong>${escapeHTML(stock.name)}</strong>
        <span>${escapeHTML(stock.symbol)} · ${escapeHTML(stock.sourceType)} · ${escapeHTML(displayText(stock.industry, "未分类"))}</span>
      </button>
    `).join("")
    : `<div class="search-empty">暂无匹配标的</div>`;
  elements.searchResults.classList.toggle("active", Boolean(elements.searchInput.value.trim()));
}

function openFirstSearchResult() {
  const match = globalSearchMatches(elements.searchInput.value)[0];
  if (!match) return false;
  window.location.hash = stockHash(match.symbol);
  elements.searchResults?.classList.remove("active");
  return true;
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

function renderDividendPanel(stock, isHolding) {
  const dividend = stock.dividend;
  if (!dividend) {
    return `
      <section class="panel dividend-panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Dividend</p>
            <h2>股息与现金流</h2>
          </div>
        </div>
        <div class="empty-state compact-empty">暂无股息数据</div>
      </section>
    `;
  }

  const currencyCode = dividendCurrency(stock);
  const perShare = finiteNumber(dividend.dividendPerShare);
  const annualLocal = dividendAnnualCashLocal(stock);
  const annualCny = dividendAnnualCashCny(stock);
  const estimatedCash = finiteNumber(dividend.estimatedAnnualCash);
  const dividendYield = calculatedDividendYield(stock);
  const shareholderReturnYield = calculatedShareholderReturnYield(stock);
  const reliability = dividendReliability(stock);
  const shield = dividendShield(stock);
  const forecastPerShare = finiteNumber(dividend.forecastPerShare);
  const forecastYield = forecastDividendYield(stock);
  const forecastCurrency = String(dividend.forecastCurrency || currencyCode).toUpperCase();

  return `
    <section class="panel dividend-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Dividend</p>
          <h2>股息与现金流</h2>
        </div>
      </div>
      <div class="detail-content">
        <div class="dividend-grid">
          <div><span>财年</span><strong>${escapeHTML(displayText(dividend.fiscalYear, "-"))}</strong></div>
          <div><span>每股分红</span><strong>${Number.isFinite(perShare) ? currency(perShare, currencyCode) : "-"}</strong></div>
          <div><span>股息率</span><strong>${displayDividendRatio(dividendYield)}</strong></div>
          <div><span>回报门槛</span><strong>${displayDividendRatio(shield.target)}</strong><small>${marketKind(stock) === "HK" ? "H股主策略" : "A股主策略"}</small></div>
          <div><span>综合回报率</span><strong>${displayDividendRatio(shareholderReturnYield)}</strong></div>
          <div><span>股息可靠性</span><strong>${badge(reliability.text, reliability.tone)}</strong></div>
          <div><span>预估财年</span><strong>${escapeHTML(displayText(dividend.forecastFiscalYear, "-"))}</strong></div>
          <div><span>预估每股</span><strong>${Number.isFinite(forecastPerShare) ? currency(forecastPerShare, forecastCurrency) : "-"}</strong></div>
          <div><span>预估股息率</span><strong>${displayDividendRatio(forecastYield)}</strong><small>参考，不作为主策略门槛</small></div>
          <div><span>预估年现金</span><strong>${Number.isFinite(annualLocal) ? currency(annualLocal, currencyCode) : Number.isFinite(estimatedCash) ? currency(estimatedCash, currencyCode) : "-"}</strong><small>${isHolding && annualCny > 0 ? `折人民币 ${currency(annualCny)}` : ""}</small></div>
        </div>
      </div>
    </section>
  `;
}

function renderNetCashPanel(stock) {
  const strategy = strategyProfile(stock);
  const netCash = strategy.netCash;
  const localCurrency = netCash.currency || stock.currency || "CNY";
  const fcfYears = Number.isFinite(netCash.fcfPositiveYears)
    ? `${netCash.fcfPositiveYears} 年`
    : Number.isFinite(netCash.fcfRecord)
      ? `${Math.round(netCash.fcfRecord * 5)} / 5 年`
      : "-";
  const haircutText = Number.isFinite(netCash.haircut) ? percent(netCash.haircut * 100, false) : "-";
  const reason = netCash.reason || (strategy.bucket === "excluded" ? "重大风险折扣归零" : "按分红稳定性与风险文本自动分档");

  return `
    <section class="panel source-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Net Cash</p>
          <h2>净现金烟蒂口径</h2>
        </div>
      </div>
      <div class="source-grid">
        <div>
          <span>策略归属</span>
          <strong>${escapeHTML(strategy.status)}</strong>
          <small>${escapeHTML(strategy.blockers.length ? strategy.blockers.join("；") : "当前无主要阻碍")}</small>
        </div>
        <div>
          <span>净现金</span>
          <strong>${financialAmount(netCash.netCash, localCurrency)}</strong>
          <small>现金/短投 - 有息债务</small>
        </div>
        <div>
          <span>净现金折扣</span>
          <strong>${haircutText}</strong>
          <small>${escapeHTML(reason)}</small>
        </div>
        <div>
          <span>调整后净现金</span>
          <strong>${financialAmount(netCash.adjustedLocal, localCurrency)}</strong>
          <small>${Number.isFinite(netCash.adjustedCny) ? `折人民币 ${financialAmount(netCash.adjustedCny, "CNY")}` : ""}</small>
        </div>
        <div>
          <span>ex-cash PE</span>
          <strong>${financialMultiple(netCash.exCashPe)}</strong>
          <small>${marketKind(stock) === "HK" ? "H股" : "A股"}门槛 ≤${strategy.peLimit}x</small>
        </div>
        <div>
          <span>ex-cash P/FCF</span>
          <strong>${financialMultiple(netCash.exCashPfcf)}</strong>
          <small>${escapeHTML(netCash.fcfBasis || "普通股东 FCF")}</small>
        </div>
        <div>
          <span>FCF yield</span>
          <strong>${financialRatio(netCash.fcfYield)}</strong>
          <small>${escapeHTML(netCash.fcfBasis || "普通股东 FCF")} / 总市值</small>
        </div>
        <div>
          <span>普通股东 FCF</span>
          <strong>${financialAmount(netCash.fcfLocal, netCash.fcfCurrency)}</strong>
          <small>${Number.isFinite(netCash.minorityFcfAdjustment) ? `少数股东分流 ${financialAmount(netCash.minorityFcfAdjustment, netCash.fcfCurrency)}` : "用于 P/FCF 和 FCF yield"}</small>
        </div>
        <div>
          <span>FCF连续性</span>
          <strong>${escapeHTML(fcfYears)}</strong>
          <small>最近五年可得数据</small>
        </div>
      </div>
    </section>
  `;
}

function renderOwnerAuditPanel(stock) {
  const audit = strategyProfile(stock).ownerAudit;
  const auditSummary = audit.blockers.length
    ? audit.blockers.join("；")
    : audit.corruptedNotes
      ? `${audit.corruptedNotes} 项说明乱码，需重建证据`
      : `评分达到 ${OWNER_AUDIT_SCORE_TARGET}/100 门槛`;
  return `
    <section class="panel source-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Owner Cash Flow</p>
          <h2>长期股东现金流评分</h2>
        </div>
      </div>
      <div class="source-grid">
        <div>
          <span>总评分</span>
          <strong>${badge(audit.text, audit.tone)}</strong>
          <small>${escapeHTML(audit.hasAudit ? `${audit.score}/100 · ${audit.grade}；${auditSummary}` : auditSummary)}</small>
        </div>
        ${audit.items.map((item) => `
          <div>
            <span>${escapeHTML(item.label)}</span>
            <strong>${badge(`${item.score}/${item.maxScore}`, item.tone)}</strong>
            <small>${escapeHTML(`${item.text}；${ownerAuditEvidenceText(item, audit.hasAudit)}`)}</small>
          </div>
        `).join("")}
      </div>
    </section>
  `;
}

function ownerAuditEvidenceText(item, hasAudit) {
  if (displayText(item.note, "")) return item.note;
  const evidence = item.evidence || "补充可验证证据";
  if (item.noteCorrupted) return `原说明乱码，需重建：${evidence}`;
  if (!hasAudit) return `未导入评分证据：${evidence}`;
  if (item.status === "pass") return `已标通过，但缺少证据说明：${evidence}`;
  if (item.status === "fail") return `已标失败，需补失败依据：${evidence}`;
  return `${item.core ? "核心项待复核" : "待复核"}：${evidence}`;
}

function renderDataSourcePanel(stock) {
  const dividend = stock.dividend ?? {};
  const hasCashDividendFormula = Number.isFinite(finiteNumber(dividend.cashDividendTotal)) && Number.isFinite(finiteNumber(stock.marketCap));
  const hasBuyback = Number.isFinite(finiteNumber(dividend.buybackAmount)) && finiteNumber(dividend.buybackAmount) > 0;
  const dividendFormula = hasCashDividendFormula
    ? "现金分红总额 / 总市值"
    : dividend.dividendPerShare
      ? "每股分红 / 最新价"
      : "暂无";
  const returnFormula = hasBuyback ? "(现金分红总额 + 回购金额) / 总市值" : "同股息率";

  return `
    <section class="panel source-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Source</p>
          <h2>数据口径</h2>
        </div>
      </div>
      <div class="source-grid">
        <div>
          <span>行情</span>
          <strong>${escapeHTML(displayText(stock.currentPriceDate, "收盘日未知"))}</strong>
          <small>最新价与昨收来自行情更新</small>
        </div>
        <div>
          <span>内在价值</span>
          <strong>${Number.isFinite(stock.intrinsicValue) ? currency(stock.intrinsicValue, stock.currency) : "-"}</strong>
          <small>由分析 JSON 导入</small>
        </div>
        <div>
          <span>安全边际</span>
          <strong>${displayMarginOfSafety(stock)}</strong>
          <small>(内在价值 - 最新价) / 内在价值</small>
        </div>
        <div>
          <span>买入分层</span>
          <strong>${displayPriceLevel(stock, "initialBuyPrice")}</strong>
          <small>旧分层保留为参考；主策略另需安全边际≥15%</small>
        </div>
        <div>
          <span>双策略门槛</span>
          <strong>A股6% / H股8%</strong>
          <small>最近完整财年综合回报率达标</small>
        </div>
        <div>
          <span>烟蒂门槛</span>
          <strong>A股≤10x / H股≤8x</strong>
          <small>按折扣后净现金扣减市值后计算 PE</small>
        </div>
        <div>
          <span>股息率</span>
          <strong>${escapeHTML(dividendFormula)}</strong>
          <small>${escapeHTML(displayText(dividend.fiscalYear, "财年未知"))}</small>
        </div>
        <div>
          <span>综合回报率</span>
          <strong>${escapeHTML(returnFormula)}</strong>
          <small>${hasBuyback ? "包含回购" : "未录入回购金额"}</small>
        </div>
      </div>
    </section>
  `;
}

function financialMetricCard(label, value, detail = "") {
  return `
    <div>
      <span>${escapeHTML(label)}</span>
      <strong>${value}</strong>
      <small>${escapeHTML(detail)}</small>
    </div>
  `;
}

function renderFinancialTable(stock, annual, currencyCode) {
  return `
    <div class="financial-table-wrap">
      <table class="financial-table">
        <thead>
          <tr>
            <th>年度</th>
            <th>收入</th>
            <th>归母利润</th>
            <th>经营现金流</th>
            <th>FCF</th>
            <th>ROE/ROIC</th>
            <th>负债率</th>
            <th>毛利/净利率</th>
            <th>PE/PB</th>
            <th>存货/应收</th>
          </tr>
        </thead>
        <tbody>
          ${annual.map((item, index) => `
            <tr>
              <td>
                <strong>${escapeHTML(displayText(item.fiscalYear, item.reportDate || "-"))}</strong>
                <small>${escapeHTML(displayText(item.reportType, item.reportDate || ""))}</small>
              </td>
              <td>
                ${financialAmount(item.revenue, item.currency || currencyCode)}
                <small>${financialRatio(item.revenueYoY, "")}</small>
              </td>
              <td>
                ${financialAmount(item.netProfit, item.currency || currencyCode)}
                <small>${financialRatio(item.netProfitYoY, "")}</small>
              </td>
              <td>${financialAmount(item.operatingCashFlow, item.currency || currencyCode)}</td>
              <td>${financialFcfCell(stock, item, index, currencyCode)}</td>
              <td>
                ${financialRatio(item.roe)}
                <small>${financialRatio(item.roic, "")}</small>
              </td>
              <td>${financialRatio(item.debtRatio)}</td>
              <td>
                ${financialRatio(item.grossMargin)}
                <small>${financialRatio(item.netMargin, "")}</small>
              </td>
              <td>
                ${financialMultiple(item.peAtCurrentPrice)}
                <small>${financialMultiple(item.pbAtCurrentPrice, "")}</small>
              </td>
              <td>
                ${Number.isFinite(finiteNumber(item.inventoryTurnoverDays)) ? `${finiteNumber(item.inventoryTurnoverDays).toFixed(0)}天` : "-"}
                <small>${Number.isFinite(finiteNumber(item.receivableTurnoverDays)) ? `${finiteNumber(item.receivableTurnoverDays).toFixed(0)}天` : ""}</small>
              </td>
            </tr>
          `).join("")}
        </tbody>
      </table>
    </div>
  `;
}

function financialFcfCell(stock, item, index, currencyCode) {
  const profile = netCashProfile(stock);
  if (index === 0 && Number.isFinite(profile.shareholderFcf)) {
    const basis = profile.fcfBasis || "普通股东 FCF";
    return `${financialAmount(profile.shareholderFcf, profile.fcfCurrency)}<small>${escapeHTML(basis)}</small>`;
  }
  return financialAmount(item.freeCashFlow, item.currency || currencyCode);
}

function renderFinancialsPanel(stock) {
  const financials = stock.financials ?? {};
  const annual = financialAnnuals(stock);
  if (!annual.length) {
    return `
      <section class="panel financials-panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Financials</p>
            <h2>多年财务数据</h2>
          </div>
        </div>
        <div class="empty-state compact-empty">暂无结构化多年财务数据，点击页面顶部「更新财务」拉取</div>
      </section>
    `;
  }

  const latest = annual[0] ?? {};
  const valuation = financialValuation(stock);
  const currencyCode = financials.currency || latest.currency || stock.currency;
  const netCash = netCashProfile(stock);
  const latestFcfDetail = Number.isFinite(netCash.shareholderFcf)
    ? `普通股东 FCF ${financialAmount(netCash.shareholderFcf, netCash.fcfCurrency)}`
    : `FCF ${financialAmount(latest.freeCashFlow, currencyCode)}`;
  const avgRoe = recentAverage(stock, "roe");
  const avgRoic = recentAverage(stock, "roic");
  const fcfRecord = positiveRecordRatio(stock, "freeCashFlow");

  return `
    <section class="panel financials-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Financials</p>
          <h2>多年财务数据</h2>
        </div>
        <div class="financials-meta">
          <span>${escapeHTML(displayText(financials.source, "数据源未知"))}</span>
          <strong>${escapeHTML(displayText(financials.updatedAt, "未记录更新时间"))}</strong>
        </div>
      </div>
      <div class="detail-content">
        <div class="financial-summary-grid">
          ${financialMetricCard("最新年收入", financialAmount(latest.revenue, currencyCode), `${displayText(latest.fiscalYear, "最新年")} · ${financialRatio(latest.revenueYoY, "同比未知")}`)}
          ${financialMetricCard("最新年利润", financialAmount(latest.netProfit, currencyCode), financialRatio(latest.netProfitYoY, "同比未知"))}
          ${financialMetricCard("经营现金流", financialAmount(latest.operatingCashFlow, currencyCode), latestFcfDetail)}
          ${financialMetricCard("ROE / ROIC", `${financialRatio(latest.roe)} / ${financialRatio(latest.roic)}`, `近年均值 ${financialRatio(avgRoe)} / ${financialRatio(avgRoic)}`)}
          ${financialMetricCard("负债率", financialRatio(latest.debtRatio), `FCF 为正 ${Number.isFinite(fcfRecord) ? `${Math.round(fcfRecord * 100)}%` : "未知"}`)}
          ${financialMetricCard("PE / PB / PEG", `${financialMultiple(valuation.pe)} / ${financialMultiple(valuation.pb)} / ${Number.isFinite(finiteNumber(valuation.peg)) ? finiteNumber(valuation.peg).toFixed(2) : "-"}`, `PE低/中/高 ${rangeText(valuation.peRange)}`)}
        </div>
        ${renderFinancialTable(stock, annual, currencyCode)}
        <p class="financial-source-note">${escapeHTML(displayText(valuation.sourceNote, "财务指标来自数据源披露口径；缺失字段会留空，不参与打分。"))}</p>
      </div>
    </section>
  `;
}

function killCriteriaItems(stock) {
  if (Array.isArray(stock?.killCriteria)) {
    return stock.killCriteria.map((item) => String(item ?? "").trim()).filter(Boolean);
  }
  if (String(stock?.killCriteria ?? "").trim()) {
    return [String(stock.killCriteria).trim()];
  }
  const candidates = [];
  const risk = String(stock?.risk ?? "").trim();
  const status = String(stock?.status ?? "").trim();
  if (risk) candidates.push(risk);
  if (/停牌|调查|内控|财报可信|治理|否决|低于75|自由现金流|扣非利润|毛利率|应收|减值/.test(status)) {
    candidates.push(status);
  }
  return candidates.slice(0, 3);
}

function renderKillCriteriaPanel(stock) {
  const items = killCriteriaItems(stock);
  return `
    <section class="panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Bear Case</p>
          <h2>我可能错在哪里</h2>
        </div>
      </div>
      <div class="detail-content">
        ${items.length
          ? `<div class="risk-check-list">${items.map((item) => `<div>${escapeHTML(item)}</div>`).join("")}</div>`
          : `<div class="empty-state compact-empty">暂无明确反证条件，建议下次分析补充</div>`}
      </div>
    </section>
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
  const confidence = confidenceMeta(stock);
  const strategy = strategyProfile(stock);

  elements.stockDetail.innerHTML = `
    <section class="detail-hero">
      <a class="ghost-button detail-back" href="#positions">返回持仓</a>
      <div>
        <p class="eyebrow">${escapeHTML(stock.symbol)} · ${escapeHTML(displayText(stock.industry, "未分类"))}</p>
        <h2>${escapeHTML(stock.name)}</h2>
      </div>
      <div class="detail-hero-actions">
        <button class="ghost-button compact-link" type="button" data-update-financials="${escapeHTML(stock.symbol)}">
          <span>↻</span>
          更新财务
        </button>
        <div class="detail-hero-meta">
          <span>${escapeHTML(closeDateText(stock) || displayText(stock.updatedAt, "行情日期未知"))}</span>
          <small>${escapeHTML(stock.financials?.updatedAt ? `财务 ${stock.financials.updatedAt}` : "财务数据待更新")}</small>
        </div>
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

    ${renderMasterVotesPanel(stock, totalValue)}

    <nav class="detail-section-nav" aria-label="详情分段导航">
      <button type="button" data-detail-section="detailValuation">估值质量</button>
      <button type="button" data-detail-section="detailOwnerAudit">长期评分</button>
      <button type="button" data-detail-section="detailFinancials">多年财务</button>
      <button type="button" data-detail-section="detailIncome">股息/净现金</button>
      <button type="button" data-detail-section="detailRisk">风险反证</button>
      <button type="button" data-detail-section="detailRecords">财报日志</button>
    </nav>

    <section class="detail-section" id="detailValuation">
      <div class="detail-section-head">
        <p class="eyebrow">Valuation</p>
        <h2>估值质量</h2>
      </div>
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
            <div><dt>策略归属</dt><dd>${escapeHTML(strategy.status)}</dd></div>
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
            <div><span>主策略安全边际</span><strong>${Number.isFinite(strategy.margin) ? percent(strategy.margin * 100, false) : "-"} / ${percent(MAIN_DCF_MARGIN_TARGET * 100, false)}</strong></div>
            <div><span>综合回报率</span><strong>${displayDividendRatio(strategy.shield.value)} / ${displayDividendRatio(strategy.shield.target)}</strong><small>${escapeHTML(strategy.shield.source)}</small></div>
            <div><span>长期评分</span><strong>${badge(`${strategy.ownerAudit.hasAudit ? strategy.ownerAudit.score : "-"} / 100`, strategy.ownerAudit.tone)}</strong><small>${escapeHTML(strategy.ownerAudit.text)}</small></div>
            <div><span>烟蒂PE</span><strong>${financialMultiple(strategy.netCash.exCashPe)} / ≤${strategy.peLimit}x</strong></div>
            <div><span>估值可信度</span><strong>${badge(confidence.text, confidence.tone)}</strong></div>
            <div><span>质量总分</span><strong>${qualityText}</strong></div>
            <div><span>健康状态</span><strong>${health ? `<span class="health-pill ${health.tone}">${health.status}</span>` : "-"}</strong></div>
            <div><span>健康评分</span><strong>${health ? `${health.score} 分` : "-"}</strong></div>
            <div><span>内在价值</span><strong>${Number.isFinite(stock.intrinsicValue) ? currency(stock.intrinsicValue, stock.currency) : "-"}</strong></div>
            <div><span>观察价</span><strong>${displayPriceLevel(stock, "watchPrice")}</strong></div>
            <div><span>首买价</span><strong>${displayPriceLevel(stock, "initialBuyPrice")}</strong></div>
            <div><span>重仓价</span><strong>${displayPriceLevel(stock, "aggressiveBuyPrice")}</strong></div>
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
    </section>

    <section class="detail-section" id="detailOwnerAudit">
      <div class="detail-section-head">
        <p class="eyebrow">Owner Cash Flow</p>
        <h2>长期股东现金流评分</h2>
      </div>
    ${renderOwnerAuditPanel(stock)}
    </section>

    <section class="detail-section" id="detailFinancials">
      <div class="detail-section-head">
        <p class="eyebrow">Financials</p>
        <h2>多年财务</h2>
      </div>
    ${renderFinancialsPanel(stock)}
    </section>

    <section class="detail-section" id="detailIncome">
      <div class="detail-section-head">
        <p class="eyebrow">Income</p>
        <h2>股息/净现金</h2>
      </div>
    ${renderDividendPanel(stock, isHolding)}

    ${renderNetCashPanel(stock)}

    ${renderDataSourcePanel(stock)}
    </section>

    <section class="detail-section" id="detailRisk">
      <div class="detail-section-head">
        <p class="eyebrow">Risk</p>
        <h2>风险反证</h2>
      </div>
    ${renderKillCriteriaPanel(stock)}

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
    </section>

    <section class="detail-section" id="detailRecords">
      <div class="detail-section-head">
        <p class="eyebrow">Records</p>
        <h2>财报日志</h2>
      </div>
    ${renderResearchUpdatesPanel(stock)}

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
    </section>
  `;
}

function render() {
  const positions = computePositions();
  renderQuoteUpdateStatus(positions);
  renderMetrics(positions);
  renderPositions(positions);
  renderDecisionArea(positions);
  renderCommitteeOverview(positions);
  renderMastersPage(positions);
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
  setQuoteUpdateStatus("正在拉取持仓和候选池最新日线收盘价及股息数据...");

  try {
    const result = await requestJSON("/api/quotes/update", { method: "POST" });
    state = result.state;
    localStorage.removeItem(STORAGE_KEY);
    syncCash();
    render();

    const skipped = result.skipped ?? [];
    if (skipped.length) {
      const preview = skipped.slice(0, 3).map((item) => `${item.symbol} ${item.error}`).join("；");
      setQuoteUpdateStatus(`已更新 ${result.updated} 个标的，股息现金流已刷新，${skipped.length} 个失败：${preview}`, "error");
    } else {
      setQuoteUpdateStatus(`已更新 ${result.updated} 个标的，股息现金流已刷新`, "success");
    }
  } finally {
    elements.updateQuotesButton.disabled = false;
    elements.updateQuotesButton.innerHTML = "<span>↻</span> 更新行情";
  }
}

async function updateFinancials(symbol, button) {
  if (!USE_BACKEND) throw new Error("需要通过 go run . 启动后端后才能更新财务数据");

  const originalHTML = button?.innerHTML;
  if (button) {
    button.disabled = true;
    button.innerHTML = "<span>↻</span> 更新中";
  }

  try {
    const result = await requestJSON(`/api/financials/update/${encodeURIComponent(symbol)}`, { method: "POST" });
    state = result.state;
    localStorage.removeItem(STORAGE_KEY);
    render();
    window.location.hash = stockHash(result.symbol);
  } finally {
    if (button) {
      button.disabled = false;
      button.innerHTML = originalHTML;
    }
  }
}

function filenameFromContentDisposition(header, fallback) {
  const raw = header ?? "";
  const utf8Match = raw.match(/filename\*=UTF-8''([^;]+)/i);
  if (utf8Match) return decodeURIComponent(utf8Match[1].replaceAll("\"", ""));

  const asciiMatch = raw.match(/filename="?([^";]+)"?/i);
  if (asciiMatch) return asciiMatch[1];

  return fallback;
}

async function exportChatGPTContext() {
  if (!USE_BACKEND) throw new Error("需要通过 go run . 启动后端后才能导出档案");

  elements.exportChatGPTButton.disabled = true;
  elements.exportChatGPTButton.innerHTML = "<span>↓</span> 导出中";
  setQuoteUpdateStatus("正在生成 ChatGPT 档案...");

  try {
    const response = await fetch("/api/chatgpt/export");
    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: "导出档案失败" }));
      throw new Error(error.error ?? "导出档案失败");
    }

    const blob = await response.blob();
    const filename = filenameFromContentDisposition(
      response.headers.get("Content-Disposition"),
      "portfolio-context.zip"
    );
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = filename;
    document.body.append(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);

    setQuoteUpdateStatus("ChatGPT 档案已导出", "success");
  } finally {
    elements.exportChatGPTButton.disabled = false;
    elements.exportChatGPTButton.innerHTML = "<span>↓</span> 导出档案";
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

document.querySelectorAll("[data-master-scroll]").forEach((button) => {
  button.addEventListener("click", () => {
    showEmbeddedMasters(button.dataset.masterScroll);
  });
});

function showEmbeddedMasters(target = "masters") {
  showPage("overview");
  requestAnimationFrame(() => {
    document.getElementById(target)?.scrollIntoView({ behavior: "smooth", block: "start" });
  });
}

function showPage(view) {
  const isStockDetail = view.startsWith("stock=");
  let nextView = isStockDetail ? "stock-detail" : pageTitles[view] ? view : "overview";
  if (!document.querySelector(`[data-page="${nextView}"]`)) {
    nextView = "overview";
  }

  if (isStockDetail) {
    renderStockDetail(computePositions(), decodeURIComponent(view.slice("stock=".length)));
  }

  document.querySelector(".nav-item.active")?.classList.remove("active");
  document.querySelector(`.nav-item[data-view="${nextView}"]`)?.classList.add("active");
  document.querySelector(".page.active")?.classList.remove("active");
  document.querySelector(`[data-page="${nextView}"]`)?.classList.add("active");
  elements.pageTitle.textContent = pageTitles[nextView];
}

function handleRoute(rawHash) {
  const view = rawHash || "overview";
  if (view === "masters") {
    showEmbeddedMasters("masters");
    return;
  }
  showPage(view);
}

window.addEventListener("hashchange", () => {
  handleRoute(window.location.hash.slice(1));
});

elements.searchInput.addEventListener("input", (event) => {
  searchTerm = event.target.value.trim().toLowerCase();
  renderSearchResults();
  render();
});

elements.searchInput.addEventListener("keydown", (event) => {
  if (event.key !== "Enter") return;
  if (openFirstSearchResult()) {
    event.preventDefault();
  }
});

elements.searchResults?.addEventListener("click", (event) => {
  const button = event.target.closest("[data-search-symbol]");
  if (!button) return;
  window.location.hash = stockHash(button.dataset.searchSymbol);
  elements.searchResults.classList.remove("active");
});

document.addEventListener("click", (event) => {
  if (event.target.closest(".search-wrap")) return;
  elements.searchResults?.classList.remove("active");
});

elements.candidateSort.addEventListener("change", (event) => {
  candidateSort = event.target.value;
  renderPlanAndCandidates();
});

elements.masterMatrix?.addEventListener("click", (event) => {
  const button = event.target.closest("[data-master-matrix-sort]");
  if (!button) return;
  masterMatrixSort = {
    key: button.dataset.masterMatrixSort,
    direction: button.dataset.nextDirection === "asc" ? "asc" : "desc"
  };
  renderMasterMatrix(computePositions());
});

elements.decisionLogFilters?.addEventListener("click", (event) => {
  const button = event.target.closest("[data-log-filter]");
  if (!button) return;
  elements.decisionLogFilters.querySelector(".active")?.classList.remove("active");
  button.classList.add("active");
  decisionLogFilter = button.dataset.logFilter;
  renderDecisionLogs();
});

elements.decisionLogToggle?.addEventListener("click", () => {
  const collapsed = elements.decisionLogPanel.classList.toggle("collapsed");
  elements.decisionLogToggle.textContent = collapsed ? "展开" : "收起";
});

elements.clearDecisionLogs?.addEventListener("click", async () => {
  try {
    await clearNonTradeDecisionLogs();
  } catch (error) {
    setQuoteUpdateStatus(error.message, "error");
  }
});

document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-detail-section]");
  if (!button) return;
  document.getElementById(button.dataset.detailSection)?.scrollIntoView({ behavior: "smooth", block: "start" });
});

document.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-update-financials]");
  if (!button) return;
  try {
    await updateFinancials(button.dataset.updateFinancials, button);
  } catch (error) {
    window.alert(error.message);
  }
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

elements.exportChatGPTButton.addEventListener("click", async () => {
  try {
    await exportChatGPTContext();
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
  handleRoute(window.location.hash.slice(1));
}

init();
