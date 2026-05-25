const STORAGE_KEY = "stock-portfolio-desk-v2";
const HOLDINGS_PRIVACY_KEY = "stock-portfolio-desk-holdings-masked";
const PORTFOLIO_RETURN_SCROLL_KEY = "stock-portfolio-return-scroll-y";

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
      action: "放入普通跟踪观察；当前不买入，等待扣非利润和自由现金流验证",
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
    { rank: 99, name: "岚图汽车", priority: "普通跟踪/低优先级", advice: "HK$4.2-4.8才接近可观察买入区；若2026H1扣非利润和自由现金流转正，可重新上修估值", discipline: "质量分低于75原则上不进入核心资产池；安全边际不足时不试仓" }
  ],
  candidates: [
    {
      symbol: "600690.SH",
      name: "海尔智家",
      status: "晴仓30跟踪",
      action: "放入晴仓30跟踪；A股暂不追，H股赔率更优",
      marginOfSafety: 0.17,
      qualityScore: 83,
      industry: "家电/全球化白电/智慧家庭",
      currency: "CNY",
      intrinsicValue: 26,
      fairValueRange: "¥24-28",
      targetBuyPrice: 19.5
    }
  ],
  industries: [
    {
      id: "coal",
      name: "煤炭",
      category: "周期资源",
      status: "待建立行业框架",
      updatedAt: "2026-05-12",
      summary: "行业研究档案占位。后续补充煤价、供需、库存、进口、政策和代表公司跟踪。",
      keywords: ["煤炭", "动力煤", "焦煤"],
      linkedSymbols: [],
      keyQuestions: ["煤价中枢是否上移或下移", "供给约束是否仍有效", "高分红可持续性是否被资本开支侵蚀"],
      metrics: [
        { name: "动力煤价格", unit: "元/吨", latestValue: null, asOf: "", source: "", comment: "预留煤价追踪" }
      ],
      notes: [
        { date: "2026-05-12", title: "行业页占位", summary: "后续按煤炭行业具体研究设计指标和内容。" }
      ]
    },
    {
      id: "oil",
      name: "石油",
      category: "周期资源",
      status: "待建立行业框架",
      updatedAt: "2026-05-12",
      summary: "行业研究档案占位。后续补充油价、资本开支、储量、桶油成本和分红纪律跟踪。",
      keywords: ["石油", "油气", "海上油气"],
      linkedSymbols: ["0883.HK"],
      keyQuestions: ["油价假设是否足够保守", "桶油成本和储量替代率是否稳定", "分红与回购是否穿越周期"],
      metrics: [
        { name: "Brent 原油", unit: "美元/桶", latestValue: null, asOf: "", source: "", comment: "预留油价追踪" }
      ],
      notes: [
        { date: "2026-05-12", title: "行业页占位", summary: "后续按石油行业具体研究设计指标和内容。" }
      ]
    },
    {
      id: "paper",
      name: "造纸",
      category: "周期制造",
      status: "待建立行业框架",
      updatedAt: "2026-05-12",
      summary: "行业研究档案占位。后续补充纸价、木浆、废纸、库存、开工率和下游需求跟踪。",
      keywords: ["造纸", "纸浆", "箱板纸", "文化纸"],
      linkedSymbols: [],
      keyQuestions: ["纸价和浆价剪刀差是否改善", "库存周期处于哪个阶段", "行业新增产能是否压制盈利"],
      metrics: [
        { name: "纸浆价格", unit: "元/吨", latestValue: null, asOf: "", source: "", comment: "预留纸价/浆价追踪" }
      ],
      notes: [
        { date: "2026-05-12", title: "行业页占位", summary: "后续按造纸行业具体研究设计指标和内容。" }
      ]
    }
  ],
  rules: [
    { dimension: "商业模式", score: 30, standard: "需求刚性、收入可重复、定价权、资本开支、行业空间" },
    { dimension: "护城河", score: 25, standard: "品牌/规模/网络效应/牌照/成本优势、份额稳定、利润率优于同行" },
    { dimension: "管理层/企业文化/治理", score: 20, standard: "长期主义、资本配置、股东回报、披露透明、少画饼" },
    { dimension: "财务质量", score: 25, standard: "ROE/ROIC、自由现金流、资产负债表、利润率、应收/存货/资本开支" }
  ]
};

const palette = ["#1aa88a", "#39aee6", "#ffbe56", "#4dae69", "#e65d6a", "#c97b14"];
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
const HK_STOCK_CONNECT_DIVIDEND_TAX_RATE = 0.2;
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
const MASTER_MATRIX_FILTERS = [
  { key: "all", label: "全部" },
  { key: "holding", label: "持仓" },
  { key: "candidate", label: "跟踪" }
];
const POSITION_CATEGORY_ORDER = ["core", "repair", "tactical"];
const POSITION_CATEGORY_META = {
  core: { key: "core", label: "核心仓", tone: "core", order: 0 },
  repair: { key: "repair", label: "修复仓", tone: "repair", order: 1 },
  tactical: { key: "tactical", label: "机动仓", tone: "tactical", order: 2 }
};
const POSITION_CATEGORY_OVERRIDES = {
  "600036.SH": "core",
  "0506.HK": "core",
  "600887.SH": "core",
  "0700.HK": "core",
  "000333.SZ": "core",
  "002415.SZ": "core",
  "0696.HK": "core",
  "600563.SH": "repair",
  "2669.HK": "repair",
  "6049.HK": "repair",
  "1405.HK": "tactical",
  "7489.HK": "tactical"
};

const STOCK_DETAIL_NAV_ITEMS = [
  { id: "detailInputs", label: "判断", desktopLabel: "人工判断" },
  { id: "detailSummary", label: "摘要", desktopLabel: "研究摘要" },
  { id: "detailValuation", label: "估值", desktopLabel: "估值证据" },
  { id: "detailFinancials", label: "财务", desktopLabel: "财务质量" },
  { id: "detailIncome", label: "现金", desktopLabel: "现金回报" },
  { id: "detailRisk", label: "风险", desktopLabel: "风险反证" },
  { id: "detailRecords", label: "日志", desktopLabel: "日志档案" }
];

const STOCK_DETAIL_VALUATION_SCENARIOS = [
  { key: "bear", label: "保守" },
  { key: "base", label: "基准" },
  { key: "bull", label: "乐观" }
];

const RESEARCH_DESK_FILTERS = [
  { value: "all", label: "全部" },
  { value: "queue", label: "队列" },
  { value: "opportunity", label: "机会" },
  { value: "data", label: "数据" },
  { value: "industry", label: "行业" }
];

let state = loadState();
let activeFilter = "all";
let positionSort = { key: "", direction: "desc" };
let sunny30Sort = { key: "quality", direction: "desc" };
let researchDeskFilter = "all";
const expandedPositionCards = new Set();
const expandedSunny30Cards = new Set();
const expandedStockDetailSections = new Set(["detailInputs", "detailSummary", "detailValuation", "detailFinancials"]);
const expandedResearchDeskSections = new Set(["queue", "opportunity"]);
let activeStockDetailSection = "detailInputs";
let holdingsMasked = localStorage.getItem(HOLDINGS_PRIVACY_KEY) === "1";
let pendingResearch = null;
let candidateSort = "consensus";
let candidateFilter = "all";
let decisionLogFilter = "all";
let masterMatrixSort = { key: "margin", direction: "desc" };
let masterMatrixFilter = "all";
let backendStateError = "";
let backendAvailable = false;
const pageTitles = {
  overview: "晴仓记",
  screener: "选股估值",
  holdings: "持仓",
  logs: "日志",
  industry: "研究台",
  trades: "日志",
  "industry-detail": "行业分析",
  "stock-detail": "股票分析详情"
};
let activeRoute = routeInfo(window.location.hash.slice(1));

const elements = {
  pageTitle: document.querySelector("#pageTitle"),
  positionCategorySummary: document.querySelector("#positionCategorySummary"),
  positionMobileSort: document.querySelector("#positionMobileSort"),
  positionMobileCards: document.querySelector("#positionMobileCards"),
  sunny30Summary: document.querySelector("#sunny30Summary"),
  screeningWeightsPanel: document.querySelector("#screeningWeightsPanel"),
  sunny30MobileSort: document.querySelector("#sunny30MobileSort"),
  sunny30MobileCards: document.querySelector("#sunny30MobileCards"),
  sunny30Body: document.querySelector("#sunny30Body"),
  positionsBody: document.querySelector("#positionsBody"),
  tradeList: document.querySelector("#tradeList"),
  allocationChart: document.querySelector("#allocationChart"),
  assetAllocationBar: document.querySelector("#assetAllocationBar"),
  assetAllocationDonut: document.querySelector("#assetAllocationDonut"),
  allocationCenterLabel: document.querySelector("#allocationCenterLabel"),
  allocationCenterValue: document.querySelector("#allocationCenterValue"),
  assetAllocationTitle: document.querySelector("#assetAllocationTitle"),
  assetAllocationSummary: document.querySelector("#assetAllocationSummary"),
  assetAllocationLegend: document.querySelector("#assetAllocationLegend"),
  valuationModuleList: document.querySelector("#valuationModuleList"),
  updateValuationHistoryButton: document.querySelector("#updateValuationHistory"),
  overviewPlanList: document.querySelector("#overviewPlanList"),
  committeeConsensus: document.querySelector("#committeeConsensus"),
  decisionQueueCount: document.querySelector("#decisionQueueCount"),
  overviewBuyCandidates: document.querySelector("#overviewBuyCandidates"),
  overviewRiskReview: document.querySelector("#overviewRiskReview"),
  disciplineDashboard: document.querySelector("#disciplineDashboard"),
  dataQualityList: document.querySelector("#dataQualityList"),
  decisionLogList: document.querySelector("#decisionLogList"),
  decisionLogPanel: document.querySelector("#decisionLogPanel"),
  decisionLogToggle: document.querySelector("#decisionLogToggle"),
  decisionLogFilters: document.querySelector("#decisionLogFilters"),
  masterMatrix: document.querySelector("#masterMatrix"),
  masterMatrixFilters: document.querySelector("#masterMatrixFilters"),
  grahamSummary: document.querySelector("#grahamSummary"),
  grahamList: document.querySelector("#grahamList"),
  buffettSummary: document.querySelector("#buffettSummary"),
  buffettList: document.querySelector("#buffettList"),
  candidateList: document.querySelector("#candidateList"),
  candidateSort: document.querySelector("#candidateSort"),
  industryList: document.querySelector("#industryList"),
  industryDetail: document.querySelector("#industryDetail"),
  stockDetail: document.querySelector("#stockDetail"),
  totalAssetsMetric: document.querySelector("#totalAssetsMetric"),
  totalValue: document.querySelector("#totalValue"),
  totalPositionPnl: document.querySelector("#totalPositionPnl"),
  totalPositionPnlRate: document.querySelector("#totalPositionPnlRate"),
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
  privacyToggle: document.querySelector("#privacyToggle"),
  updateQuotesButton: document.querySelector("#updateQuotesButton"),
  quoteUpdateStatus: document.querySelector("#quoteUpdateStatus"),
  tradeDialog: document.querySelector("#tradeDialog"),
  tradeForm: document.querySelector("#tradeForm"),
  tradeStockNames: document.querySelector("#tradeStockNames"),
  tradeNameLabel: document.querySelector("#tradeNameLabel"),
  tradePriceLabel: document.querySelector("#tradePriceLabel"),
  tradeSharesLabel: document.querySelector("#tradeSharesLabel"),
  tradeNameInput: document.querySelector("#tradeNameInput"),
  tradeSharesInput: document.querySelector("#tradeSharesInput"),
  openSunny30CandidateButton: document.querySelector("#openSunny30Candidate"),
  sunny30CandidateDialog: document.querySelector("#sunny30CandidateDialog"),
  sunny30CandidateForm: document.querySelector("#sunny30CandidateForm"),
  holdingDialog: document.querySelector("#holdingDialog"),
  holdingForm: document.querySelector("#holdingForm"),
  researchDialog: document.querySelector("#researchDialog"),
  researchForm: document.querySelector("#researchForm"),
  researchJSON: document.querySelector("#researchJSON"),
  researchPreview: document.querySelector("#researchPreview"),
  researchStatus: document.querySelector("#researchStatus"),
  importResearchButton: document.querySelector("#importResearch"),
  backToTopButton: document.querySelector("#backToTopButton")
};

const HOLDINGS_MASK = "******";

function privateText(value, mask = HOLDINGS_MASK) {
  return holdingsMasked ? mask : value;
}

function privateHTML(value, mask = HOLDINGS_MASK) {
  return holdingsMasked ? escapeHTML(mask) : value;
}

function privateClass(className) {
  return holdingsMasked ? "" : className;
}

function setPrivacyToggleState() {
  if (!elements.privacyToggle) return;
  const label = holdingsMasked ? "显示持仓数据" : "隐藏持仓数据";
  elements.privacyToggle.classList.toggle("is-masked", holdingsMasked);
  elements.privacyToggle.setAttribute("aria-pressed", String(holdingsMasked));
  elements.privacyToggle.setAttribute("aria-label", label);
  elements.privacyToggle.setAttribute("title", label);
  document.body.classList.toggle("holdings-masked", holdingsMasked);
}

if (!USE_BACKEND) {
  syncCash();
}

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
  const { headers = {}, timeoutMs = 0, signal, ...fetchOptions } = options;
  const timeoutController = timeoutMs > 0 && !signal ? new AbortController() : null;
  const timeout = timeoutController
    ? window.setTimeout(() => timeoutController.abort(), timeoutMs)
    : null;

  const response = await fetch(path, {
    cache: "no-store",
    ...fetchOptions,
    signal: signal ?? timeoutController?.signal,
    headers: { "Content-Type": "application/json", ...headers }
  }).catch((error) => {
    if (error?.name === "AbortError") {
      throw new Error(`${path} 请求超时`);
    }
    throw error;
  }).finally(() => {
    if (timeout) window.clearTimeout(timeout);
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: "request failed" }));
    throw new Error(error.error ?? "request failed");
  }

  return response.json();
}

function normalizedLoadedState(rawState) {
  const stocks = Array.isArray(rawState?.stocks) ? rawState.stocks : [];
  const holdings = stocks.length ? stocksToHoldings(stocks) : (Array.isArray(rawState?.holdings) ? rawState.holdings : []);
  const candidates = stocks.length ? stocksToCandidates(stocks) : (Array.isArray(rawState?.candidates) ? rawState.candidates : []);
  return {
    ...structuredClone(seedState),
    ...(rawState ?? {}),
    stocks,
    holdings,
    candidates,
    trades: Array.isArray(rawState?.trades) ? rawState.trades : [],
    decisionLogs: Array.isArray(rawState?.decisionLogs) ? rawState.decisionLogs : [],
    screeningWeights: rawState?.screeningWeights ?? { quality: 30, cashFlow: 25, valuation: 20, shareholderReturn: 15, growth: 10 },
    plan: Array.isArray(rawState?.plan) ? rawState.plan : [],
    industries: Array.isArray(rawState?.industries) ? rawState.industries : [],
    rules: Array.isArray(rawState?.rules) ? rawState.rules : []
  };
}

function setLoadedState(rawState) {
  state = normalizedLoadedState(rawState);
  return state;
}

function stocksToHoldings(stocks) {
  return stocks
    .filter((stock) => stock?.position)
    .map((stock) => ({
      ...stock,
      shares: Number(stock.position?.shares) || 0,
      cost: Number(stock.position?.cost) || 0
    }));
}

function stocksToCandidates(stocks) {
  return stocks.map((stock) => {
    const { position, ...candidate } = stock;
    return candidate;
  });
}

function stockPayloadFromLegacy(stock) {
  const payload = { ...(stock ?? {}) };
  const shares = finiteNumber(payload.shares);
  const cost = finiteNumber(payload.cost);
  if (Number.isFinite(shares) || Number.isFinite(cost)) {
    payload.position = {
      shares: Number.isFinite(shares) ? shares : 0,
      cost: Number.isFinite(cost) ? cost : 0
    };
  }
  delete payload.shares;
  delete payload.cost;
  return payload;
}

function applyStaticRuntimeQuote(stock, record) {
  if (!stock || !record) return stock;
  const next = { ...stock };
  const currentPrice = finiteNumber(record.currentPrice);
  if (Number.isFinite(currentPrice) && currentPrice > 0) {
    next.currentPrice = currentPrice;
    next.marginOfSafety = calculatedMarginOfSafety({ ...next, currentPrice }) ?? next.marginOfSafety;
  }
  ["previousClose", "twentyDayClose"].forEach((key) => {
    const value = finiteNumber(record[key]);
    if (Number.isFinite(value) && value > 0) next[key] = value;
  });
  ["twentyDayCloseDate", "currentPriceDate", "previousCloseDate", "updatedAt"].forEach((key) => {
    const value = String(record[key] ?? "").trim();
    if (value) next[key] = value;
  });
  if (record.twentyDayChange !== undefined && record.twentyDayChange !== null) {
    const value = finiteNumber(record.twentyDayChange);
    if (Number.isFinite(value)) next.twentyDayChange = value;
  }
  const marketCap = finiteNumber(record.marketCap);
  if (Number.isFinite(marketCap) && marketCap > 0) {
    next.marketCap = marketCap;
    next.marketCapCurrency = String(record.marketCapCurrency || record.currency || next.currency || "").trim().toUpperCase();
  }
  if (!String(next.currency ?? "").trim() && record.currency) {
    next.currency = String(record.currency).trim().toUpperCase();
  }
  const dividendPerShare = finiteNumber(record.dividendPerShare);
  if (Number.isFinite(dividendPerShare) && dividendPerShare > 0) {
    next.dividend = { ...(next.dividend ?? {}), dividendPerShare };
    if (record.dividendCurrency) next.dividend.dividendCurrency = String(record.dividendCurrency).trim().toUpperCase();
    if (record.dividendFiscalYear) next.dividend.fiscalYear = String(record.dividendFiscalYear).trim();
  }
  return next;
}

function mergeStaticRuntimeQuotes(nextState, quoteBook) {
  const quotes = quoteBook?.quotes ?? {};
  const quoteFor = (symbol) => quotes[normalizeSymbol(symbol)];
  return {
    ...nextState,
    holdings: (nextState.holdings ?? []).map((holding) => applyStaticRuntimeQuote(holding, quoteFor(holding.symbol))),
    candidates: (nextState.candidates ?? []).map((candidate) => applyStaticRuntimeQuote(candidate, quoteFor(candidate.symbol)))
  };
}

function staticDataStatus(apiError, quoteBook) {
  const quoteCount = Object.keys(quoteBook?.quotes ?? {}).length;
  const issues = [{
    tone: "info",
    title: "静态数据模式",
    detail: `/api/state 不可用，已从 data/portfolio.json 加载。${apiError?.message ? `原因：${apiError.message}` : ""}`
  }];
  if (!quoteCount) {
    issues.push({
      tone: "warn",
      title: "行情缓存缺失",
      detail: "未读取到 data/runtime/quotes.json，行情会使用 portfolio.json 内的旧价格。"
    });
  }
  return {
    status: quoteCount ? "ok" : "warn",
    dataDir: "data",
    writable: false,
    backupCount: 0,
    portfolio: { path: "data/portfolio.json", exists: true },
    runtimeQuotes: {
      path: "data/runtime/quotes.json",
      exists: quoteCount > 0,
      updatedAt: String(quoteBook?.updatedAt ?? "")
    },
    runtimeIndustryMetrics: { path: "data/runtime/industry_metrics.json", exists: false },
    issues
  };
}

async function loadStaticState(apiError) {
  const staticState = normalizedLoadedState(await requestJSON("./data/portfolio.json", { timeoutMs: 5000 }));
  const quoteBook = await requestJSON("./data/runtime/quotes.json", { timeoutMs: 5000 }).catch(() => null);
  state = mergeStaticRuntimeQuotes(staticState, quoteBook);
  state.dataStatus = staticDataStatus(apiError, quoteBook);
  backendStateError = "";
  backendAvailable = false;
  localStorage.removeItem(STORAGE_KEY);
  return true;
}

async function loadBackendState() {
  if (!USE_BACKEND) return false;

  try {
    setLoadedState(await requestJSON("/api/state", { timeoutMs: 3000 }));
    backendStateError = "";
    backendAvailable = true;
    localStorage.removeItem(STORAGE_KEY);
    return true;
  } catch (error) {
    console.warn("后端不可用，尝试加载静态数据", error);
    try {
      return await loadStaticState(error);
    } catch (staticError) {
      console.warn("静态数据也不可用，使用浏览器本地数据", staticError);
      backendAvailable = false;
      backendStateError = `${error.message || "后端不可用"}；静态数据加载失败：${staticError.message || staticError}`;
      setQuoteUpdateStatus("后端和静态数据都不可用，已切换到浏览器本地兜底数据", "error");
      return false;
    }
  }
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

function industryHash(id) {
  return `#industry=${encodeURIComponent(id)}`;
}

function normalizeSymbol(symbol) {
  const text = String(symbol ?? "").trim().toUpperCase();
  if (/^HK\d+$/.test(text)) {
    const code = text.slice(2);
    return `${String(Number(code)).padStart(4, "0")}.HK`;
  }
  if (text.endsWith(".HK")) {
    const code = text.slice(0, -3);
    if (/^\d+$/.test(code)) {
      return `${String(Number(code)).padStart(4, "0")}.HK`;
    }
  }
  return text;
}

function normalizeAssetType(assetType) {
  return String(assetType ?? "").trim().toLowerCase() === "fund" ? "fund" : "stock";
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
  const currentPrice = Number.isFinite(position?.currentPrice) && position.currentPrice > 0
    ? currency(position.currentPrice, position.currency)
    : "价格未知";
  const previousClose = Number.isFinite(position?.previousClose) && position.previousClose > 0
    ? currency(position.previousClose, position.currency)
    : "价格未知";
  const currentDate = position.currentPriceDate || "未知";
  const previousDate = position.previousCloseDate || "未知";
  return `今收 ${currentPrice} · ${currentDate}；昨收 ${previousClose} · ${previousDate}`;
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

function normalizeRate(value, fallback = 0) {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return fallback;
  if (number > 1 && number <= 100) return clamp(number / 100, 0, 1);
  return clamp(number, 0, 1);
}

function dividendTaxProfile(stock) {
  const dividend = stock?.dividend ?? {};
  const stockConnectRate = marketKind(stock) === "HK"
    ? normalizeRate(
        dividend.stockConnectDividendTaxRate ??
        dividend.stockConnectTaxRate ??
        dividend.personalDividendTaxRate,
        HK_STOCK_CONNECT_DIVIDEND_TAX_RATE
      )
    : 0;
  const nonResidentRate = normalizeRate(
    dividend.nonResidentWithholdingTaxRate ??
    dividend.foreignWithholdingTaxRate ??
    dividend.withholdingTaxRate,
    0
  );
  const creditable = Boolean(dividend.withholdingTaxCreditable);
  const factor = creditable
    ? 1 - Math.max(stockConnectRate, nonResidentRate)
    : (1 - nonResidentRate) * (1 - stockConnectRate);
  return {
    stockConnectRate,
    nonResidentRate,
    creditable,
    factor: clamp(factor, 0, 1),
    effectiveTaxRate: 1 - clamp(factor, 0, 1)
  };
}

function dividendTaxRate(stock) {
  return dividendTaxProfile(stock).effectiveTaxRate;
}

function dividendTaxFactor(stock) {
  return dividendTaxProfile(stock).factor;
}

function dividendTaxText(stock) {
  const profile = dividendTaxProfile(stock);
  const note = displayText(stock?.dividend?.taxNote, "");
  if (profile.stockConnectRate <= 0 && profile.nonResidentRate <= 0) return note || "未折扣分红税";
  const parts = [];
  if (profile.stockConnectRate > 0) parts.push(`港股通个税 ${percent(profile.stockConnectRate * 100, false)}`);
  if (profile.nonResidentRate > 0) parts.push(`非居民预提 ${percent(profile.nonResidentRate * 100, false)}`);
  parts.push(`${profile.creditable ? "抵免后" : "到账"}税负 ${percent(profile.effectiveTaxRate * 100, false)}`);
  if (note) parts.push(note);
  return parts.join(" · ");
}

function afterTaxDividendCny(stock, grossCny) {
  const amount = finiteNumber(grossCny);
  if (!Number.isFinite(amount)) return null;
  return amount * dividendTaxFactor(stock);
}

function dividendYieldInputs(stock) {
  const dividend = stock?.dividend;
  const cashDividendTotal = finiteNumber(dividend?.cashDividendTotal);
  const marketCap = finiteNumber(stock?.marketCap);
  if (Number.isFinite(cashDividendTotal) && cashDividendTotal > 0 && Number.isFinite(marketCap) && marketCap > 0) {
    const grossCashDividendTotalCny = cashDividendTotal * fx(cashDividendTotalCurrency(stock));
    return {
      cashDividendTotalCny: afterTaxDividendCny(stock, grossCashDividendTotalCny),
      grossCashDividendTotalCny,
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
    cashDividendTotalCny: afterTaxDividendCny(stock, perShare * fx(dividendCurrency(stock))),
    grossCashDividendTotalCny: perShare * fx(dividendCurrency(stock)),
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
  return afterTaxDividendCny(stock, localCash * fx(dividendCurrency(stock))) ?? 0;
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
  if (Number.isFinite(explicit) && explicit > 0) return explicit * dividendTaxFactor(stock);
  const perShare = finiteNumber(dividend?.forecastPerShare);
  const currentPrice = finiteNumber(stock?.currentPrice);
  if (!Number.isFinite(perShare) || perShare <= 0 || !Number.isFinite(currentPrice) || currentPrice <= 0) return null;
  const forecastCurrency = String(dividend?.forecastCurrency || dividendCurrency(stock)).toUpperCase();
  return (perShare * fx(forecastCurrency) * dividendTaxFactor(stock)) / (currentPrice * fx(stock?.currency || forecastCurrency));
}

function dividendShield(stock) {
  const trailing = calculatedDividendYield(stock);
  const shareholderReturn = calculatedShareholderReturnYield(stock);
  const forecast = forecastDividendYield(stock);
  const value = Number.isFinite(shareholderReturn) ? shareholderReturn : null;
  const source = Number.isFinite(shareholderReturn)
    ? marketKind(stock) === "HK" ? "最近财年综合回报（港股通股息税后）" : "最近财年综合回报"
    : "未记录综合回报";
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
  const netCashCny = toCnyAmount(netCash, profileCurrency);
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
    netCashCny,
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
  if (!mainPassed && !cigarPassed && bucket !== "excluded" && Number.isFinite(netCash.exCashPe)) {
    blockers.push(`烟蒂PE需≤${peLimit}x且FCF达标`);
  }

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

function cleanDecisionReason(value) {
  return String(value ?? "")
    .trim()
    .replace(/^已触发重大风险否决项[:：]\s*/, "")
    .replace(/^未达标[（(]\s*/, "")
    .replace(/[）)]$/, "")
    .trim();
}

function statusReasonText(status) {
  const text = displayText(status, "");
  const match = text.match(/[（(]([^（）()]+)[）)]/);
  return cleanDecisionReason(match?.[1] || text);
}

function riskReasonText(stock) {
  const reason = cleanDecisionReason(stock?.risk);
  return reason && !/^无/.test(reason) ? reason : "";
}

function noteCoreReasonText(stock) {
  const text = displayText(stock?.notes, "");
  const match = text.match(/核心判断[:：]([^。；]+)/);
  return cleanDecisionReason(match?.[1] || "");
}

function stockActionReasonText(stock, strategy) {
  const isHolding = stock?.sourceType === "holding";
  const statusReason = statusReasonText(stock?.status);
  const actionText = displayText(stock?.action, "");

  if (strategy.bucket === "excluded") {
    const reason = riskReasonText(stock) || statusReason || noteCoreReasonText(stock) || actionText || "风险补偿不足";
    return `${isHolding ? "卖出理由" : "不买入理由"}：${reason}`;
  }

  if (strategy.bucket === "main") {
    const reasons = [
      Number.isFinite(strategy.shield.value) ? `${strategy.shield.source} ${displayDividendRatio(strategy.shield.value)}达标` : "",
      Number.isFinite(strategy.margin) ? `安全边际 ${percent(strategy.margin * 100, false)}` : "",
      strategy.ownerAudit.hasAudit ? `长期股东评分 ${strategy.ownerAudit.score}/100` : ""
    ].filter(Boolean);
    return `${isHolding ? "加仓理由" : "买入理由"}：${reasons.length ? reasons.join("；") : statusReason || actionText || "主策略条件达标"}`;
  }

  if (strategy.bucket === "cigar") {
    const reasons = [
      `调整后净现金 ${financialAmount(strategy.netCash.adjustedCny, "CNY")}`,
      `ex-cash PE ${financialMultiple(strategy.netCash.exCashPe)}`,
      `FCF ${financialMultiple(strategy.netCash.exCashPfcf)}`
    ].filter((item) => !item.endsWith(" -"));
    return `${isHolding ? "加仓理由" : "买入理由"}：${reasons.length ? reasons.join("；") : statusReason || actionText || "辅策略烟蒂条件达标"}`;
  }

  const waitReason = statusReason || actionText || strategy.blockers.slice(0, 2).join("；") || "关键买入条件未完全达标";
  return `等待理由：${waitReason}`;
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
  if (backendStateError) {
    setQuoteUpdateStatus(`后端不可用：${backendStateError}`, "error");
    return;
  }
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
  if (type === "candidate") return "更新跟踪标的";
  if (type === "newCandidate") return "新增到晴仓30";
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
  const realizedPnl = realizedStockPnlBySymbol();
  return state.holdings
    .filter((holding) => holding.shares > 0)
    .map((holding) => {
      const symbol = normalizeSymbol(holding.symbol);
      const shares = finiteNumber(holding.shares) ?? 0;
      const cost = finiteNumber(holding.cost) ?? 0;
      const currentPrice = finiteNumber(holding.currentPrice);
      const previousClose = finiteNumber(holding.previousClose);
      const hasCurrentPrice = Number.isFinite(currentPrice) && currentPrice > 0;
      const marketValueLocal = hasCurrentPrice ? shares * currentPrice : 0;
      const costValueLocal = shares * cost;
      const marketValueCny = marketValueLocal * fx(holding.currency);
      const costValueCny = costValueLocal * fx(holding.currency);
      const unrealizedPnlCny = hasCurrentPrice ? marketValueCny - costValueCny : null;
      const realized = realizedPnl.get(symbol) || { pnlCny: 0, costCny: 0 };
      const realizedPnlCny = realized.pnlCny;
      const realizedCostCny = realized.costCny;
      const pnlCny = Number.isFinite(unrealizedPnlCny)
        ? unrealizedPnlCny + realizedPnlCny
        : Math.abs(realizedPnlCny) > 0 ? realizedPnlCny : null;
      const pnlCostValueCny = costValueCny + realizedCostCny;
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
        realizedCostCny,
        pnlCostValueCny,
        unrealizedPnlCny,
        realizedPnlCny,
        pnlCny,
        pnlRate: pnlCostValueCny && Number.isFinite(pnlCny) ? (pnlCny / pnlCostValueCny) * 100 : null,
        dayChange
      };
    });
}

function assetSummary(positions = computePositions()) {
  const stockValue = positions.reduce((sum, item) => sum + (finiteNumber(item.marketValueCny) ?? 0), 0);
  const cashValue = finiteNumber(state.cash) ?? 0;
  const totalAssets = stockValue + cashValue;
  return { stockValue, cashValue, totalAssets };
}

function sortedStockTrades() {
  return [...(state.trades ?? [])]
    .map((trade, index) => ({ ...trade, index }))
    .filter((trade) => normalizeAssetType(trade.assetType) !== "fund" && normalizeSymbol(trade.symbol))
    .sort((a, b) => {
      const idA = finiteNumber(a.id);
      const idB = finiteNumber(b.id);
      if (Number.isFinite(idA) && Number.isFinite(idB) && idA !== idB) return idA - idB;
      const dateCompare = String(a.date ?? "").localeCompare(String(b.date ?? ""));
      return dateCompare || a.index - b.index;
    });
}

function realizedStockPnlBySymbol() {
  const trades = sortedStockTrades();
  const openingLots = new Map();

  state.holdings.forEach((holding) => {
    const symbol = normalizeSymbol(holding.symbol);
    const shares = finiteNumber(holding.shares) ?? 0;
    const cost = finiteNumber(holding.cost) ?? 0;
    if (!symbol || shares <= 0) return;
    openingLots.set(symbol, {
      shares,
      cost,
      currency: String(holding.currency || "CNY").toUpperCase()
    });
  });

  [...trades].reverse().forEach((trade) => {
    const symbol = normalizeSymbol(trade.symbol);
    const shares = finiteNumber(trade.shares) ?? 0;
    const price = finiteNumber(trade.price);
    const side = String(trade.side ?? "").trim().toLowerCase();
    if (!symbol || shares <= 0 || !Number.isFinite(price)) return;
    const lot = openingLots.get(symbol) || {
      shares: 0,
      cost: price,
      currency: String(trade.currency || "CNY").toUpperCase()
    };

    if (side === "buy") {
      const previousShares = lot.shares - shares;
      if (previousShares > 0) {
        const previousCostTotal = lot.shares * lot.cost - shares * price;
        lot.shares = previousShares;
        if (Number.isFinite(previousCostTotal)) lot.cost = previousCostTotal / previousShares;
      } else {
        lot.shares = 0;
        lot.cost = price;
      }
      openingLots.set(symbol, lot);
      return;
    }

    if (side === "sell") {
      lot.shares += shares;
      openingLots.set(symbol, lot);
    }
  });

  const lots = new Map([...openingLots.entries()].map(([symbol, lot]) => [symbol, { ...lot }]));
  const realized = new Map();

  trades.forEach((trade) => {
    const symbol = normalizeSymbol(trade.symbol);
    const shares = finiteNumber(trade.shares) ?? 0;
    const price = finiteNumber(trade.price);
    const side = String(trade.side ?? "").trim().toLowerCase();
    if (!symbol || shares <= 0 || !Number.isFinite(price)) return;
    const currencyCode = String(trade.currency || lots.get(symbol)?.currency || "CNY").toUpperCase();
    const lot = lots.get(symbol) || { shares: 0, cost: price, currency: currencyCode };

    if (side === "buy") {
      const totalCost = lot.shares * lot.cost + shares * price;
      lot.shares += shares;
      if (lot.shares > 0) lot.cost = totalCost / lot.shares;
      lot.currency = currencyCode;
      lots.set(symbol, lot);
      return;
    }

    if (side === "sell") {
      const multiplier = fx(currencyCode);
      const entry = realized.get(symbol) || { pnlCny: 0, costCny: 0 };
      entry.pnlCny += shares * (price - lot.cost) * multiplier;
      entry.costCny += shares * lot.cost * multiplier;
      realized.set(symbol, entry);
      lot.shares = Math.max(0, lot.shares - shares);
      lot.currency = currencyCode;
      lots.set(symbol, lot);
    }
  });

  return realized;
}

function syncCash() {
  if (Number.isFinite(state.cash)) return;

  const investedStocks = state.holdings.reduce((sum, holding) => {
    const shares = finiteNumber(holding.shares) ?? 0;
    const currentPrice = finiteNumber(holding.currentPrice) ?? 0;
    return sum + shares * currentPrice * fx(holding.currency);
  }, 0);
  state.cash = state.totalCapital - investedStocks;
}

function getFilteredPositions(positions) {
  return positions
    .map((stock) => ({ ...stock, sourceType: "holding" }))
    .map((stock) => ({ stock, strategy: strategyProfile(stock) }))
    .filter(({ stock }) => {
      const matchesFilter =
        activeFilter === "all" ||
        stock.sourceType === activeFilter;

      return matchesFilter;
    });
}

function normalizePositionCategory(value) {
  const text = String(value ?? "").trim().toLowerCase();
  if (!text) return "";
  if (["core", "核心仓"].includes(text) || text.includes("核心仓")) return "core";
  if (["repair", "修复仓"].includes(text) || text.includes("修复仓")) return "repair";
  if (["tactical", "机动仓"].includes(text) || text.includes("机动仓")) return "tactical";
  return "";
}

function positionCategory(stock, strategy = strategyProfile(stock)) {
  const explicit = [
    stock?.positionCategory,
    stock?.holdingCategory,
    stock?.stockCategory,
    stock?.classification,
    stock?.category
  ].map(normalizePositionCategory).find(Boolean);
  if (explicit) return POSITION_CATEGORY_META[explicit];

  const override = POSITION_CATEGORY_OVERRIDES[normalizeSymbol(stock?.symbol)];
  if (override) return POSITION_CATEGORY_META[override];

  const text = [stock?.symbol, stock?.name, stock?.industry, stock?.status, stock?.action, stock?.notes].filter(Boolean).join(" ");
  const dividendYield = calculatedDividendYield(stock);
  const latest = latestAnnualFinancial(stock);
  const revenueGrowth = finiteNumber(latest?.revenueYoY);
  const profitGrowth = finiteNumber(latest?.netProfitYoY);
  const coreLike = /核心|复利|长期|龙头|现金流|分红|股息|银行|家电|乳制品|饮料|白酒|腾讯|美的|伊利|招商银行|中国食品|茅台|海尔/.test(text);
  const tacticalLike = /机动|套利|事件|景气|周期|小仓|弹性|油气|能源|金铜|资源|汽车|新能源|餐饮|连锁|QSR|达势|岚图|泡泡玛特/.test(text);
  const repairLike = /修复|低估|低预期|物业|地产链|安防|民航|广告|分众|同仁堂|中海物业|保利物业|海康|中航信/.test(text);

  if (tacticalLike || (Number.isFinite(revenueGrowth) && revenueGrowth >= 0.18) || (Number.isFinite(profitGrowth) && profitGrowth >= 0.22)) {
    return POSITION_CATEGORY_META.tactical;
  }
  if (repairLike || strategy.bucket === "cigar") {
    return POSITION_CATEGORY_META.repair;
  }
  if (coreLike || (Number.isFinite(dividendYield) && dividendYield >= 0.04) || strategy.bucket === "main") {
    return POSITION_CATEGORY_META.core;
  }
  return POSITION_CATEGORY_META.repair;
}

function positionCategoryPill(category) {
  const label = category?.label ?? category?.text ?? "-";
  const tone = category?.tone ?? "";
  return `<span class="position-category-pill ${tone}">${escapeHTML(label)}</span>`;
}

function positionCategorySelect(stock, category = positionCategory(stock)) {
  const symbol = normalizeSymbol(stock?.symbol);
  const options = POSITION_CATEGORY_ORDER.map((key) => {
    const item = POSITION_CATEGORY_META[key];
    const selected = item.key === category.key ? " selected" : "";
    return `<option value="${escapeHTML(item.label)}"${selected}>${escapeHTML(item.label)}</option>`;
  }).join("");
  return `
    <select class="position-category-select ${category.tone}" data-stock-category-select data-symbol="${escapeHTML(symbol)}" data-current-category="${escapeHTML(category.label)}" aria-label="选择${escapeHTML(stock?.name || symbol)}分类">
      ${options}
    </select>
  `;
}

function stockDetailCategoryControl(stock) {
  const category = positionCategory(stock);
  return `
    <label class="detail-category-control">
      <span>分类</span>
      ${positionCategorySelect(stock, category)}
    </label>
  `;
}

function renderPositionCategorySummary(positions) {
  if (!elements.positionCategorySummary) return;
  const totalValue = positions.reduce((sum, stock) => sum + (finiteNumber(stock.marketValueCny) ?? 0), 0);
  const summary = POSITION_CATEGORY_ORDER.reduce((result, key) => {
    result[key] = { value: 0, count: 0 };
    return result;
  }, {});

  positions.forEach((stock) => {
    const category = positionCategory(stock);
    summary[category.key].value += finiteNumber(stock.marketValueCny) ?? 0;
    summary[category.key].count += 1;
  });

  elements.positionCategorySummary.innerHTML = POSITION_CATEGORY_ORDER.map((key) => {
    const category = POSITION_CATEGORY_META[key];
    const item = summary[key];
    const share = totalValue > 0 ? item.value / totalValue : 0;
    return `
      <div class="position-category-summary-cell ${category.tone}">
        <span>${escapeHTML(category.label)}</span>
        <strong>${escapeHTML(privateText(percent(share * 100, false)))}</strong>
        <small>${escapeHTML(privateText(`${item.count} 只 · ${wholeCurrency(item.value)}`))}</small>
      </div>
    `;
  }).join("");
}

function positionSortValue(item, key) {
  const { stock, strategy } = item;
  if (key === "category") return positionCategory(stock, strategy).order;
  if (key === "marketValue") return finiteNumber(stock.marketValueCny);
  if (key === "pnl") return finiteNumber(stock.pnlCny);
  if (key === "return") return finiteNumber(strategy?.shield?.value);
  if (key === "owner") return strategy?.ownerAudit?.hasAudit ? finiteNumber(strategy.ownerAudit.score) : null;
  if (key === "margin") return finiteNumber(strategy?.margin);
  return null;
}

function compareNullableNumbers(a, b, direction = "desc") {
  const aNumber = finiteNumber(a);
  const bNumber = finiteNumber(b);
  const aValid = Number.isFinite(aNumber);
  const bValid = Number.isFinite(bNumber);
  if (!aValid && !bValid) return 0;
  if (!aValid) return 1;
  if (!bValid) return -1;
  return direction === "asc" ? aNumber - bNumber : bNumber - aNumber;
}

function sortedPositions(positions) {
  const filtered = getFilteredPositions(positions);
  if (!positionSort.key) return filtered;
  return [...filtered].sort((a, b) => {
    const result = compareNullableNumbers(positionSortValue(a, positionSort.key), positionSortValue(b, positionSort.key), positionSort.direction);
    return result || String(a.stock?.name ?? "").localeCompare(String(b.stock?.name ?? ""), "zh-CN");
  });
}

function positionSortDefaultDirection(key) {
  return key === "category" ? "asc" : "desc";
}

function positionMobileSortValue() {
  return positionSort.key ? `${positionSort.key}:${positionSort.direction}` : "default";
}

function updatePositionSortControls() {
  document.querySelectorAll("[data-position-sort]").forEach((button) => {
    const key = button.dataset.positionSort;
    const active = positionSort.key === key;
    const directionText = positionSort.direction === "asc" ? "升序" : "降序";
    const label = button.querySelector("span")?.textContent?.trim() || "当前列";
    button.classList.toggle("active", active);
    button.dataset.nextDirection = active && positionSort.direction === "desc" ? "asc" : "desc";
    button.setAttribute("aria-label", active ? `${label}${directionText}` : `按${label}排序`);
    button.closest("th")?.setAttribute("aria-sort", active ? (positionSort.direction === "asc" ? "ascending" : "descending") : "none");
    const indicator = button.querySelector("strong");
    if (indicator) {
      indicator.textContent = active ? (positionSort.direction === "asc" ? "↑" : "↓") : "↕";
    }
  });
  if (elements.positionMobileSort) {
    elements.positionMobileSort.value = positionMobileSortValue();
  }
}

function renderMobileStat(label, value, className = "", detail = "") {
  return `
    <div class="mobile-card-stat">
      <span>${escapeHTML(label)}</span>
      <strong class="${className || ""}">${escapeHTML(value)}</strong>
      ${detail ? `<small>${escapeHTML(detail)}</small>` : ""}
    </div>
  `;
}

function renderMobileDetail(label, value, detail = "", className = "") {
  return `
    <div class="mobile-card-detail-item">
      <span>${escapeHTML(label)}</span>
      <strong class="${className || ""}">${escapeHTML(value)}</strong>
      ${detail ? `<small>${escapeHTML(detail)}</small>` : ""}
    </div>
  `;
}

function renderPositionMobileCards(items, totalValue) {
  if (!elements.positionMobileCards) return;
  const mobileItems = items;

  if (!mobileItems.length) {
    elements.positionMobileCards.innerHTML = `<div class="empty-state compact-empty">暂无符合筛选的持仓</div>`;
    return;
  }

  elements.positionMobileCards.innerHTML = mobileItems.map(({ stock, strategy }) => {
    const symbol = normalizeSymbol(stock.symbol);
    const expanded = expandedPositionCards.has(symbol);
    const category = positionCategory(stock, strategy);
    const health = decisionHealthCell(stock, totalValue);
    const marginTone = decisionMarginTone(strategy.margin);
    const returnTone = strategy.shield.passed ? "core" : "reduce";
    const ownerTone = decisionToneClass(strategy.ownerAudit.tone);
    const dayClass = privateClass(stock.dayChange >= 0 ? "positive" : "negative");
    const pnlClass = privateClass(stock.pnlCny >= 0 ? "positive" : "negative");
    const dayRate = stock.marketValueCny ? percent((stock.dayChange / stock.marketValueCny) * 100) : "";
    const weight = totalValue ? (stock.marketValueCny / totalValue) * 100 : 0;
    const currentPrice = Number.isFinite(stock.currentPrice) && stock.currentPrice > 0
      ? currency(stock.currentPrice, stock.currency)
      : "-";
    const quoteDate = stock.currentPriceDate || "收盘日未知";
    const sourceClass = marketKind(stock) === "HK" ? "hk" : "";
    const ownerText = strategy.ownerAudit.hasAudit ? `${strategy.ownerAudit.score}/100` : "待评分";
    return `
      <article class="mobile-position-card ${expanded ? "is-expanded" : ""}">
        <div class="mobile-card-toggle" aria-expanded="${expanded ? "true" : "false"}">
          <button class="ticker mobile-card-code-action ${sourceClass}" type="button" data-toggle-position-card="${escapeHTML(symbol)}" aria-label="${expanded ? "收起" : "展开"}${escapeHTML(stock.name)}卡片">
            ${escapeHTML(stock.symbol.slice(0, 4))}
          </button>
          <a class="mobile-card-title mobile-card-title-link" href="${stockHash(stock.symbol)}" aria-label="查看${escapeHTML(stock.name)}详情">
            <strong>${escapeHTML(stock.name)}</strong>
            <small>${escapeHTML(stock.symbol)} · ${escapeHTML(currentPrice)} · ${escapeHTML(quoteDate)}</small>
          </a>
          <button class="mobile-card-pills mobile-card-inline-toggle" type="button" data-toggle-position-card="${escapeHTML(symbol)}" aria-label="${expanded ? "收起" : "展开"}${escapeHTML(stock.name)}卡片">
            ${positionCategoryPill(category)}
            <span class="health-status-score ${health.tone}">${escapeHTML(health.scoreText)}</span>
          </button>
          <button class="mobile-card-chevron mobile-card-inline-toggle" type="button" data-toggle-position-card="${escapeHTML(symbol)}">${expanded ? "收起" : "展开"}</button>
        </div>
        <div class="mobile-card-core">
          ${renderMobileStat("今日盈亏", privateText(currency(stock.dayChange)), dayClass, privateText(dayRate))}
          ${renderMobileStat("安全边际", Number.isFinite(strategy.margin) ? percent(strategy.margin * 100, false) : "-", `health-pill ${marginTone}`)}
          ${renderMobileStat("累计盈亏", privateText(currency(stock.pnlCny)), pnlClass, privateText(percent(stock.pnlRate)))}
        </div>
        <div class="mobile-card-expanded">
          <div class="mobile-card-detail-grid">
            ${renderMobileDetail("市值/现价", privateText(currency(stock.marketValueCny)), `${currentPrice} · ${quoteDate}`)}
            ${renderMobileDetail("持股/成本", privateText(`${stock.shares} 股`), `成本 ${privateText(currency(stock.cost, stock.currency))}`)}
            ${renderMobileDetail("综合回报率", privateText(`${displayDividendRatio(strategy.shield.value)} / ${displayDividendRatio(strategy.shield.target)}`), "", `health-pill ${returnTone}`)}
            ${renderMobileDetail("长期评分", ownerText, "", `health-pill ${ownerTone}`)}
            ${renderMobileDetail("仓位", privateText(`${weight.toFixed(1)}%`), privateText(health.title))}
          </div>
          <a class="ghost-button compact-link mobile-card-link" href="${stockHash(stock.symbol)}">查看详情</a>
        </div>
      </article>
    `;
  }).join("");
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

  if (highRisk || (Number.isFinite(quality) && quality < 70)) {
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
  const holdingsValue = totalValue;
  const totalCost = positions.reduce((sum, item) => {
    return sum + (finiteNumber(item.pnlCostValueCny) ?? finiteNumber(item.costValueCny) ?? 0);
  }, 0);
  const totalPositionPnl = positions.reduce((sum, item) => sum + (finiteNumber(item.pnlCny) ?? 0), 0);
  const totalPositionPnlRate = totalCost ? (totalPositionPnl / totalCost) * 100 : 0;
  const cashValue = finiteNumber(state.cash) ?? 0;
  const totalAssets = holdingsValue + cashValue;
  const dayChange = positions.reduce((sum, item) => sum + item.dayChange, 0);
  const dividends = dividendSummary(positions);
  const dividendYield = totalValue ? dividends.annualCashCny / totalValue : 0;

  elements.totalAssetsMetric.textContent = privateText(wholeCurrency(totalAssets));
  elements.totalValue.textContent = privateText(wholeCurrency(holdingsValue));
  if (elements.totalPositionPnl) {
    elements.totalPositionPnl.textContent = privateText(wholeCurrency(totalPositionPnl));
    elements.totalPositionPnl.className = privateClass(totalPositionPnl >= 0 ? "positive" : "negative");
  }
  if (elements.totalPositionPnlRate) {
    elements.totalPositionPnlRate.textContent = privateText(`相对成本 ${percent(totalPositionPnlRate)}`);
    elements.totalPositionPnlRate.className = privateClass(totalPositionPnl >= 0 ? "positive" : "negative");
  }
  if (elements.dayChange) {
    elements.dayChange.textContent = privateText(wholeCurrency(dayChange));
    elements.dayChange.className = privateClass(dayChange >= 0 ? "positive" : "negative");
  }
  if (elements.dayChangeRate) {
    elements.dayChangeRate.textContent = privateText(percent(totalValue ? (dayChange / totalValue) * 100 : 0));
    elements.dayChangeRate.className = privateClass(dayChange >= 0 ? "positive" : "negative");
  }
  elements.annualDividend.textContent = privateText(wholeCurrency(dividends.annualCashCny));
  elements.portfolioDividendYield.textContent = privateText(dividends.topContributor
    ? `组合税后股息率 ${percent(dividendYield * 100, false)} · 高风险 ${percent(dividends.annualCashCny ? (dividends.highRiskCashCny / dividends.annualCashCny) * 100 : 0, false)}`
    : "组合税后股息率 0.00%");
  elements.positionCount.textContent = privateText(`${positions.length} 只股票`);
  elements.recordCount.textContent = privateText(`${state.stocks?.length ?? state.holdings.length + state.candidates.length} 只股票 · ${state.trades.length} 条交易`);
}

function renderAssetAllocation(positions) {
  if ((!elements.assetAllocationDonut && !elements.assetAllocationBar) || !elements.assetAllocationLegend) return;
  const { stockValue, cashValue, totalAssets } = assetSummary(positions);
  const chartCash = Math.max(cashValue, 0);
  const chartTotal = stockValue + chartCash;
  const circumference = 2 * Math.PI * 72;
  const segments = [
    { key: "stock", label: "股票", value: stockValue, chartValue: stockValue, color: "#087f5b" },
    { key: "cash", label: "现金", value: cashValue, chartValue: chartCash, color: "#d69b2d" }
  ];
  if (elements.assetAllocationDonut) {
    let offset = 0;
    elements.assetAllocationDonut.innerHTML = !holdingsMasked && chartTotal > 0
      ? segments.map((segment) => {
          const length = (segment.chartValue / chartTotal) * circumference;
          const html = `<circle cx="95" cy="95" r="72" stroke="${segment.color}" stroke-dasharray="${length} ${circumference}" stroke-dashoffset="${-offset}"></circle>`;
          offset += length;
          return html;
        }).join("")
      : "";
  }

  if (elements.assetAllocationBar) {
    elements.assetAllocationBar.innerHTML = !holdingsMasked && chartTotal > 0
      ? segments.map((segment) => {
          const width = (segment.chartValue / chartTotal) * 100;
          return `<span style="width: ${Math.max(0, width).toFixed(2)}%; background: ${segment.color}" title="${escapeHTML(segment.label)}"></span>`;
        }).join("")
      : "";
  }

  const stockShare = totalAssets > 0 ? (stockValue / totalAssets) * 100 : 0;
  if (elements.allocationCenterLabel) elements.allocationCenterLabel.textContent = "股票占比";
  if (elements.allocationCenterValue) elements.allocationCenterValue.textContent = privateText(`${stockShare.toFixed(1)}%`);
  if (elements.assetAllocationTitle) elements.assetAllocationTitle.textContent = "股票 / 现金配置";
  if (elements.assetAllocationSummary) elements.assetAllocationSummary.textContent = "按当前股票市值和现金余额拆分总资产；行业暴露和单票风险在持仓页复盘。";
  const exposureRows = segments.map((segment) => {
    const share = totalAssets > 0 ? (segment.value / totalAssets) * 100 : 0;
    const negativeNote = segment.key === "cash" && segment.value < 0 ? " · 现金为负" : "";
    return `
      <div class="overview-exposure-row asset-allocation-legend-item">
        <span class="allocation-dot" style="background: ${segment.color}"></span>
        <span>${escapeHTML(segment.label)}</span>
        <span>${escapeHTML(privateText(`${wholeCurrency(segment.value)}${negativeNote}`))}</span>
        <strong>${escapeHTML(privateText(`${share.toFixed(1)}%`))}</strong>
      </div>
    `;
  });
  const topIndustryText = topGroups(positions, (item) => firstIndustry(item.industry) || "未分类", (item) => item.marketValueCny)
    .map((item) => item.name)
    .join(" / ") || "-";
  exposureRows.push(`
    <div class="overview-exposure-row industry-row">
      <span>前三行业</span>
      <strong>${escapeHTML(topIndustryText)}</strong>
    </div>
  `);
  elements.assetAllocationLegend.innerHTML = exposureRows.join("");
}

function decisionToneClass(tone) {
  if (tone === "strong" || tone === "safe" || tone === "buy") return "core";
  if (tone === "risk" || tone === "reduce") return "risk";
  if (tone === "watch" || tone === "wait") return "reduce";
  return "watch";
}

function decisionMarginTone(value) {
  if (!Number.isFinite(value)) return "watch";
  if (value >= MAIN_DCF_MARGIN_TARGET) return "core";
  if (value < 0.1) return "risk";
  return "reduce";
}

function decisionStockMeta(stock) {
  if (stock.sourceType === "holding") {
    return `${stock.symbol} · ${privateText(`${stock.shares} 股`)} · 成本 ${privateText(currency(stock.cost, stock.currency))}`;
  }
  return `${stock.symbol} · 跟踪 · ${displayText(firstIndustry(stock.industry), "未分类")}`;
}

function decisionMarketCell(stock) {
  const currentPrice = finiteNumber(stock.currentPrice);
  const priceText = Number.isFinite(currentPrice) && currentPrice > 0
    ? currency(currentPrice, stock.currency)
    : "-";
  const dateText = stock.currentPriceDate || "收盘日未知";
  if (stock.sourceType !== "holding") {
    return `<strong>-</strong><br /><small class="quote-date">${priceText} · ${escapeHTML(dateText)}</small>`;
  }
  return `<strong>${escapeHTML(privateText(currency(stock.marketValueCny)))}</strong><br /><small class="quote-date">${priceText} · ${escapeHTML(dateText)}</small>`;
}

function decisionPnlCell(stock) {
  if (stock.sourceType !== "holding") {
    return { className: "", html: `<strong>-</strong><br /><small>未持仓</small>` };
  }
  const hasRealizedPnl = Math.abs(finiteNumber(stock.realizedPnlCny) ?? 0) >= 0.005;
  const detailText = hasRealizedPnl && Number.isFinite(stock.unrealizedPnlCny)
    ? `浮动 ${currency(stock.unrealizedPnlCny)} · 已实现 ${currency(stock.realizedPnlCny)}`
    : percent(stock.pnlRate);
  return {
    className: privateClass(stock.pnlCny >= 0 ? "positive" : "negative"),
    html: `<span class="decision-primary">${escapeHTML(privateText(currency(stock.pnlCny)))}</span><br /><small>${escapeHTML(privateText(detailText))}</small>`
  };
}

function decisionHealthCell(stock, totalValue) {
  if (stock.sourceType !== "holding") {
    const qualityText = Number.isFinite(stock.qualityScore) ? `质量${stock.qualityScore}` : "质量待补";
    return {
      tone: "watch",
      scoreText: qualityText,
      weightText: "未持仓",
      title: displayText(stock.status, "跟踪标的")
    };
  }
  const health = holdingHealth(stock, totalValue);
  const weight = totalValue ? (stock.marketValueCny / totalValue) * 100 : 0;
  return {
    tone: health.tone,
    scoreText: `${health.score}分`,
    weightText: privateText(`仓位${weight.toFixed(1)}%`),
    title: health.detail
  };
}

function sunny30Universe(positions) {
  const holdingsBySymbol = new Map(positions.map((position) => [normalizeSymbol(position.symbol), position]));
  const seen = new Set();
  const stocks = (state.candidates ?? [])
    .map((candidate) => {
      const symbol = normalizeSymbol(candidate.symbol);
      if (!symbol || seen.has(symbol)) return null;
      seen.add(symbol);
      const holding = holdingsBySymbol.get(symbol);
      const merged = mergeSunny30Stock(holding, candidate, symbol);
      return {
        ...merged,
        currentPrice: finiteNumber(merged.currentPrice),
        previousClose: finiteNumber(merged.previousClose),
        marginOfSafety: calculatedMarginOfSafety(merged) ?? merged.marginOfSafety,
        strategy: strategyProfile(merged)
      };
    })
    .filter((stock) => stock?.symbol && stock?.name)
    .sort((a, b) => {
      const qualityA = finiteNumber(a.qualityScore) ?? 0;
      const qualityB = finiteNumber(b.qualityScore) ?? 0;
      const ownerA = a.strategy.ownerAudit.hasAudit ? a.strategy.ownerAudit.score : 0;
      const ownerB = b.strategy.ownerAudit.hasAudit ? b.strategy.ownerAudit.score : 0;
      return qualityB - qualityA || ownerB - ownerA || a.name.localeCompare(b.name, "zh-CN");
    });
  return sortedSunny30Stocks(stocks);
}

function mergeSunny30Stock(holding, candidate, symbol) {
  if (!holding) {
    return { ...candidate, symbol, sourceType: "candidate" };
  }

  const merged = { ...holding, symbol, sourceType: "holding" };
  ["name", "industry", "status", "action", "currency", "notes"].forEach((key) => {
    const value = String(candidate?.[key] ?? "").trim();
    if (value) merged[key] = candidate[key];
  });
  ["currentPriceDate", "previousCloseDate", "twentyDayCloseDate", "updatedAt"].forEach((key) => {
    const value = String(candidate?.[key] ?? "").trim();
    if (value) merged[key] = candidate[key];
  });
  [
    "currentPrice",
    "previousClose",
    "twentyDayClose",
    "twentyDayChange",
    "intrinsicValue",
    "targetBuyPrice",
    "marginOfSafety",
    "qualityScore",
    "businessModel",
    "moat",
    "governance",
    "financialQuality"
  ].forEach((key) => {
    const value = finiteNumber(candidate?.[key]);
    if (value !== null) merged[key] = value;
  });
  return merged;
}

function sunny30Type(stock) {
  return positionCategory(stock);
}

function sunny30Quality(stock) {
  const score = finiteNumber(stock?.qualityScore);
  if (Number.isFinite(score)) {
    if (score >= 85) return { text: "优秀", tone: "strong" };
    if (score >= 75) return { text: "合格", tone: "watch" };
    return { text: "待验证", tone: "risk" };
  }
  const audit = ownerAuditProfile(stock);
  if (audit.hasAudit && audit.score >= 85) return { text: "优秀", tone: "strong" };
  if (audit.hasAudit && audit.score >= OWNER_AUDIT_SCORE_TARGET) return { text: "合格", tone: "watch" };
  return { text: "待验证", tone: "watch" };
}

function sunny30Moat(stock) {
  const moat = finiteNumber(stock?.moat);
  if (Number.isFinite(moat)) {
    if (moat >= 23) return { text: "极强", tone: "strong" };
    if (moat >= 20) return { text: "强", tone: "strong" };
    if (moat >= 17) return { text: "稳固", tone: "watch" };
    return { text: "形成中", tone: "watch" };
  }
  const text = [stock?.industry, stock?.notes].filter(Boolean).join(" ");
  if (/品牌|平台|网络效应|牌照|规模|渠道/.test(text)) return { text: "待量化", tone: "watch" };
  return { text: "待补充", tone: "watch" };
}

function sunny30Ratio(current, base) {
  const currentValue = finiteNumber(current);
  const baseValue = finiteNumber(base);
  if (!Number.isFinite(currentValue) || !Number.isFinite(baseValue) || baseValue <= 0) return null;
  return (currentValue - baseValue) / baseValue;
}

function sunny30DisplayRatio(value) {
  const ratio = finiteNumber(value);
  if (!Number.isFinite(ratio)) return { text: "-", className: "muted" };
  return {
    text: percent(ratio * 100),
    className: ratio >= 0 ? "positive" : "negative"
  };
}

function sunny30DayChange(stock) {
  return sunny30DisplayRatio(sunny30Ratio(stock?.currentPrice, stock?.previousClose));
}

function sunny30TwentyDayChange(stock) {
  const explicit = finiteNumber(stock?.twentyDayChange);
  if (Number.isFinite(explicit)) return sunny30DisplayRatio(explicit);
  return sunny30DisplayRatio(sunny30Ratio(stock?.currentPrice, stock?.twentyDayClose));
}

function sunny30Margin(stock) {
  const margin = marginValue(stock);
  return {
    value: margin,
    text: Number.isFinite(margin) ? percent(margin * 100, false) : "-",
    tone: decisionMarginTone(margin)
  };
}

function sunny30Pill(value, className) {
  return `<span class="sunny30-pill ${className}">${escapeHTML(value)}</span>`;
}

function sunny30QualityValue(stock) {
  const score = finiteNumber(stock?.qualityScore);
  if (Number.isFinite(score)) return score;
  const audit = ownerAuditProfile(stock);
  return audit.hasAudit ? audit.score : null;
}

function sunny30TwentyDayRatio(stock) {
  const explicit = finiteNumber(stock?.twentyDayChange);
  if (Number.isFinite(explicit)) return explicit;
  return sunny30Ratio(stock?.currentPrice, stock?.twentyDayClose);
}

function sunny30ReturnValue(stock) {
  const strategy = stock?.strategy ?? strategyProfile(stock);
  return finiteNumber(strategy?.shield?.value);
}

function sunny30ReturnCell(stock) {
  const strategy = stock?.strategy ?? strategyProfile(stock);
  const value = finiteNumber(strategy?.shield?.value);
  const target = finiteNumber(strategy?.shield?.target);
  const tone = Number.isFinite(value) && Number.isFinite(target) && value >= target ? "core" : "reduce";
  return `<span class="health-pill ${tone}">${displayDividendRatio(value)} / ${displayDividendRatio(target)}</span>`;
}

function sunny30CanDelete(stock) {
  const positionShares = finiteNumber(stock?.position?.shares);
  const legacyShares = finiteNumber(stock?.shares);
  return stock?.sourceType !== "holding" && !(Number.isFinite(positionShares) && positionShares > 0) && !(Number.isFinite(legacyShares) && legacyShares > 0);
}

function sunny30DeleteControl(stock, mobile = false) {
  if (!sunny30CanDelete(stock)) {
    return `<span class="health-pill core ${mobile ? "mobile-card-link" : ""}">持仓中</span>`;
  }
  const mobileClass = mobile ? "ghost-button compact-link danger mobile-card-link " : "";
  const label = mobile ? "删除标的" : "删除";
  return `<button class="${mobileClass}sunny30-delete-button" type="button" data-delete-sunny30="${escapeHTML(stock.symbol)}">${label}</button>`;
}

const SCREENING_WEIGHT_FIELDS = [
  { key: "quality", label: "质量" },
  { key: "cashFlow", label: "现金流" },
  { key: "valuation", label: "估值" },
  { key: "shareholderReturn", label: "股东回报" },
  { key: "growth", label: "成长" }
];

function screeningWeights() {
  const defaults = { quality: 30, cashFlow: 25, valuation: 20, shareholderReturn: 15, growth: 10 };
  return SCREENING_WEIGHT_FIELDS.reduce((weights, field) => {
    const value = finiteNumber(state.screeningWeights?.[field.key]);
    weights[field.key] = Number.isFinite(value) ? value : defaults[field.key];
    return weights;
  }, {});
}

function screeningHardRejects(stock) {
  const rejects = [];
  const latest = latestAnnualFinancial(stock);
  const debtRatio = finiteNumber(latest.debtRatio);
  const freeCashFlow = finiteNumber(latest.freeCashFlow);
  const fcfRecord = positiveRecordRatio(stock, "freeCashFlow");
  const confidence = valuationConfidence(stock);
  if (hasMajorRisk(stock)) rejects.push("重大风险待排除");
  if (confidence === "low") rejects.push("财报/估值可信度低");
  if (Number.isFinite(fcfRecord) && fcfRecord < 0.5 && (!Number.isFinite(freeCashFlow) || freeCashFlow <= 0)) {
    rejects.push("长期 FCF 未验证");
  }
  if (!isFinancialBusiness(stock) && Number.isFinite(debtRatio) && debtRatio > 0.75) rejects.push("杠杆过高");
  return rejects;
}

function screeningSubscores(stock) {
  const latest = latestAnnualFinancial(stock);
  const valuation = financialValuation(stock);
  const margin = marginValue(stock);
  const revenueGrowth = finiteNumber(latest.revenueYoY);
  const profitGrowth = finiteNumber(latest.netProfitYoY);
  const revenueCagr = compoundGrowth(stock, "revenue");
  const profitCagr = compoundGrowth(stock, "netProfit");
  const fcfRecord = positiveRecordRatio(stock, "freeCashFlow");
  const cashConversion = finiteNumber(latest.operatingCashFlowToRevenue);
  const pePercentile = finiteNumber(valuation.pePercentile);
  const pbPercentile = finiteNumber(valuation.pbPercentile);
  const dividendYield = calculatedDividendYield(stock);

  const valuationBase = (
    scoreBand(margin, 0, 0.25, 0.35) * 0.55 +
    inverseScoreBand(pePercentile, 0.25, 0.8, 0.55) * 0.25 +
    inverseScoreBand(pbPercentile, 0.25, 0.8, 0.55) * 0.2
  );
  return {
    quality: Math.round(qualityComposite(stock) * 100),
    cashFlow: Math.round((
      (Number.isFinite(fcfRecord) ? fcfRecord : 0.45) * 0.55 +
      scoreBand(cashConversion, 0.05, 0.25, 0.45) * 0.25 +
      balanceSheetScore(stock, 0.55) * 0.2
    ) * 100),
    valuation: Math.round(valuationBase * 100),
    shareholderReturn: Math.round((
      dividendScore(stock) * 0.65 +
      scoreBand(dividendYield, 0.02, 0.07, 0.35) * 0.35
    ) * 100),
    growth: Math.round((
      scoreBand(revenueGrowth, -0.03, 0.12, 0.45) * 0.3 +
      scoreBand(profitGrowth, -0.05, 0.15, 0.45) * 0.25 +
      scoreBand(revenueCagr, 0, 0.1, 0.45) * 0.2 +
      scoreBand(profitCagr, 0, 0.12, 0.45) * 0.25
    ) * 100)
  };
}

function screeningProfile(stock) {
  const weights = screeningWeights();
  const subscores = screeningSubscores(stock);
  const rejects = screeningHardRejects(stock);
  const totalWeight = SCREENING_WEIGHT_FIELDS.reduce((sum, field) => sum + weights[field.key], 0) || 100;
  const score = Math.round(SCREENING_WEIGHT_FIELDS.reduce((sum, field) => {
    return sum + (subscores[field.key] ?? 0) * weights[field.key] / totalWeight;
  }, 0));
  const weightedScore = rejects.length ? Math.min(score, 59) : score;
  return {
    pass: rejects.length === 0,
    rejects,
    subscores,
    score: weightedScore
  };
}

function renderScreeningWeightsPanel(stocks) {
  if (!elements.screeningWeightsPanel) return;
  const weights = screeningWeights();
  const total = SCREENING_WEIGHT_FIELDS.reduce((sum, field) => sum + weights[field.key], 0);
  const profiles = stocks.map((stock) => screeningProfile(stock));
  const rejectedCount = profiles.filter((profile) => !profile.pass).length;
  elements.screeningWeightsPanel.innerHTML = `
    <form class="screening-weights-form" data-screening-weights-form>
      <div class="screening-rule-copy">
        <strong>先硬否决，再排序</strong>
        <span>默认口径：质量 30、现金流 25、估值 20、股东回报 15、成长 10。总和必须为 100。</span>
      </div>
      <div class="screening-weight-grid">
        ${SCREENING_WEIGHT_FIELDS.map((field) => `
          <label>
            <span>${escapeHTML(field.label)}</span>
            <input name="${escapeHTML(field.key)}" type="number" min="0" max="100" step="1" value="${escapeHTML(String(weights[field.key]))}" />
          </label>
        `).join("")}
      </div>
      <div class="screening-weight-actions">
        <span class="${total === 100 ? "positive" : "negative"}">合计 ${escapeHTML(String(total))}</span>
        <span>硬否决 ${escapeHTML(String(rejectedCount))} 只</span>
        <button class="ghost-button compact-link" type="submit">保存权重</button>
      </div>
    </form>
  `;
}

function screeningScoreCell(profile) {
  const tone = profile.pass && profile.score >= 80 ? "core" : profile.pass && profile.score >= 65 ? "watch" : "risk";
  return `<span class="health-pill ${tone}">${escapeHTML(String(profile.score))}/100</span>`;
}

function screeningSubscoresCell(profile) {
  return `
    <div class="screening-subscore-list">
      ${SCREENING_WEIGHT_FIELDS.map((field) => `
        <span><em>${escapeHTML(field.label)}</em><strong>${escapeHTML(String(profile.subscores[field.key] ?? "-"))}</strong></span>
      `).join("")}
    </div>
  `;
}

function sunny30SortValue(stock, key) {
  if (key === "name") return String(stock?.name ?? "");
  if (key === "type") return sunny30Type(stock).order;
  if (key === "screening") return screeningProfile(stock).score;
  if (key === "return") return sunny30ReturnValue(stock);
  if (key === "quality") return sunny30QualityValue(stock);
  if (key === "moat") return finiteNumber(stock?.moat);
  if (key === "margin") return marginValue(stock);
  if (key === "dayChange") return sunny30Ratio(stock?.currentPrice, stock?.previousClose);
  if (key === "twentyDayChange") return sunny30TwentyDayRatio(stock);
  return null;
}

function compareNullableStrings(a, b, direction = "asc") {
  const aText = String(a ?? "").trim();
  const bText = String(b ?? "").trim();
  if (!aText && !bText) return 0;
  if (!aText) return 1;
  if (!bText) return -1;
  const result = aText.localeCompare(bText, "zh-CN");
  return direction === "asc" ? result : -result;
}

function sortedSunny30Stocks(stocks) {
  return [...stocks].sort((a, b) => {
    const key = sunny30Sort.key || "quality";
    const valueA = sunny30SortValue(a, key);
    const valueB = sunny30SortValue(b, key);
    const result = key === "name"
      ? compareNullableStrings(valueA, valueB, sunny30Sort.direction)
      : compareNullableNumbers(valueA, valueB, sunny30Sort.direction);
    return result || compareNullableStrings(a.name, b.name, "asc");
  });
}

function sunny30DefaultDirection(key) {
  return key === "name" || key === "type" ? "asc" : "desc";
}

function sunny30MobileSortValue() {
  return `${sunny30Sort.key || "quality"}:${sunny30Sort.direction || "desc"}`;
}

function parseMobileSortValue(value) {
  if (value === "default") return { key: "", direction: "desc" };
  const [key, direction] = String(value ?? "").split(":");
  return {
    key: key || "",
    direction: direction === "asc" ? "asc" : "desc"
  };
}

function updateSunny30SortControls() {
  document.querySelectorAll("[data-sunny30-sort]").forEach((button) => {
    const key = button.dataset.sunny30Sort;
    const active = sunny30Sort.key === key;
    const directionText = sunny30Sort.direction === "asc" ? "升序" : "降序";
    const label = button.querySelector("span")?.textContent?.trim() || "当前列";
    button.classList.toggle("active", active);
    button.dataset.nextDirection = active && sunny30Sort.direction === "desc" ? "asc" : "desc";
    button.setAttribute("aria-label", active ? `${label}${directionText}` : `按${label}排序`);
    button.closest("th")?.setAttribute("aria-sort", active ? (sunny30Sort.direction === "asc" ? "ascending" : "descending") : "none");
    const indicator = button.querySelector("strong");
    if (indicator) {
      indicator.textContent = active ? (sunny30Sort.direction === "asc" ? "↑" : "↓") : "↕";
    }
  });
  if (elements.sunny30MobileSort) {
    elements.sunny30MobileSort.value = sunny30MobileSortValue();
  }
}

function renderSunny30MobileCards(stocks) {
  if (!elements.sunny30MobileCards) return;

  if (!stocks.length) {
    elements.sunny30MobileCards.innerHTML = `<div class="empty-state compact-empty">暂无晴仓30标的</div>`;
    return;
  }

  elements.sunny30MobileCards.innerHTML = stocks.map((stock) => {
    const symbol = normalizeSymbol(stock.symbol);
    const expanded = expandedSunny30Cards.has(symbol);
    const quality = sunny30Quality(stock);
    const margin = sunny30Margin(stock);
    const screening = screeningProfile(stock);
    const dayChange = sunny30DayChange(stock);
    const twentyDayChange = sunny30TwentyDayChange(stock);
    const sourceClass = marketKind(stock) === "HK" ? "hk" : "";
    const rejectedNameClass = screening.pass ? "" : " screening-reject-name";
    return `
      <article class="mobile-position-card mobile-sunny30-card ${expanded ? "is-expanded" : ""}">
        <div class="mobile-card-toggle" aria-expanded="${expanded ? "true" : "false"}">
          <button class="ticker mobile-card-code-action ${sourceClass}" type="button" data-toggle-sunny30-card="${escapeHTML(symbol)}" aria-label="${expanded ? "收起" : "展开"}${escapeHTML(stock.name)}卡片">
            ${escapeHTML(stock.symbol.slice(0, 4))}
          </button>
          <a class="mobile-card-title mobile-card-title-link${rejectedNameClass}" href="${stockHash(stock.symbol)}" aria-label="查看${escapeHTML(stock.name)}详情">
            <strong>${escapeHTML(stock.name)}</strong>
            <small>${escapeHTML(stock.symbol)} · ${escapeHTML(stock.currentPriceDate || "行情日期未知")}</small>
          </a>
          <button class="mobile-card-pills mobile-card-inline-toggle" type="button" data-toggle-sunny30-card="${escapeHTML(symbol)}" aria-label="${expanded ? "收起" : "展开"}${escapeHTML(stock.name)}卡片">
            ${screeningScoreCell(screening)}
          </button>
          <button class="mobile-card-chevron mobile-card-inline-toggle" type="button" data-toggle-sunny30-card="${escapeHTML(symbol)}">${expanded ? "收起" : "展开"}</button>
        </div>
        <div class="mobile-card-core">
          ${renderMobileStat("选股评分", `${screening.score}/100`, `sunny30-pill ${screening.pass ? "strong" : "risk"}`)}
          ${renderMobileStat("今日涨跌", privateText(dayChange.text), privateClass(dayChange.className))}
          ${renderMobileStat("公司质量", quality.text, `sunny30-pill ${quality.tone}`)}
        </div>
        <div class="mobile-card-expanded">
          <div class="mobile-card-detail-grid">
            ${renderMobileDetail("硬否决", screening.pass ? "通过" : screening.rejects.join(" / "))}
            ${renderMobileDetail("子分", SCREENING_WEIGHT_FIELDS.map((field) => `${field.label}${screening.subscores[field.key] ?? "-"}`).join(" · "))}
            ${renderMobileDetail("20日涨跌", privateText(twentyDayChange.text), "", privateClass(twentyDayChange.className))}
            ${renderMobileDetail("当前价格", Number.isFinite(stock.currentPrice) ? currency(stock.currentPrice, stock.currency) : "-", stock.currentPriceDate || "")}
          </div>
          ${sunny30DeleteControl(stock, true)}
        </div>
      </article>
    `;
  }).join("");
}

function renderSunny30(positions) {
  if (!elements.sunny30Body) return;
  const stocks = sunny30Universe(positions);
  updateSunny30SortControls();
  const profiles = stocks.map((stock) => screeningProfile(stock));
  const passedCount = profiles.filter((profile) => profile.pass).length;
  const buyableCount = stocks.filter((stock, index) => {
    const margin = marginValue(stock);
    return profiles[index].pass && Number.isFinite(margin) && margin >= MAIN_DCF_MARGIN_TARGET;
  }).length;
  const averageScore = profiles.length
    ? Math.round(profiles.reduce((sum, profile) => sum + profile.score, 0) / profiles.length)
    : 0;

  if (elements.sunny30Summary) {
    elements.sunny30Summary.innerHTML = [
      ["股票池", `${stocks.length}/50`],
      ["通过硬否决", `${passedCount} 只`],
      ["安全边际达标", `${buyableCount} 只`],
      ["平均选股分", `${averageScore}`]
    ].map(([label, value]) => `
      <div class="sunny30-summary-cell">
        <span>${escapeHTML(label)}</span>
        <strong>${escapeHTML(value)}</strong>
      </div>
    `).join("");
  }
  renderScreeningWeightsPanel(stocks);

  if (!stocks.length) {
    elements.sunny30Body.innerHTML = `<tr><td colspan="7" class="empty-state">暂无股票池标的</td></tr>`;
    renderSunny30MobileCards(stocks);
    return;
  }

  renderSunny30MobileCards(stocks);

  elements.sunny30Body.innerHTML = stocks.map((stock) => {
    const screening = screeningProfile(stock);
    const margin = sunny30Margin(stock);
    const dayChange = sunny30DayChange(stock);
    const twentyDayChange = sunny30TwentyDayChange(stock);
    const sourceClass = marketKind(stock) === "HK" ? "hk" : "";
    const rejectedNameClass = screening.pass ? "" : " screening-reject-name";
    const rejectedTitle = screening.pass ? "" : ` title="硬否决：${escapeHTML(screening.rejects.join(" / ") || "否决")}"`;
    const actionCell = sunny30DeleteControl(stock);
    return `
      <tr>
        <td>
          <div class="stock-cell">
            <span class="ticker ${sourceClass}">${escapeHTML(stock.symbol.slice(0, 4))}</span>
            <a class="stock-name stock-link${rejectedNameClass}" href="${stockHash(stock.symbol)}"${rejectedTitle}>
              <strong>${escapeHTML(stock.name)}</strong>
              <span>${escapeHTML(stock.symbol)}</span>
            </a>
          </div>
        </td>
        <td data-label="选股评分">${screeningScoreCell(screening)}</td>
        <td data-label="五项子分">${screeningSubscoresCell(screening)}</td>
        <td data-label="安全边际">${sunny30Pill(margin.text, margin.tone)}</td>
        <td data-label="今日涨跌"><span class="${dayChange.className}">${escapeHTML(dayChange.text)}</span></td>
        <td data-label="20日涨跌"><span class="${twentyDayChange.className}">${escapeHTML(twentyDayChange.text)}</span></td>
        <td class="sunny30-action-cell" data-label="操作">${actionCell}</td>
      </tr>
    `;
  }).join("");
}

function renderPositions(positions) {
  const filtered = sortedPositions(positions);
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  updatePositionSortControls();
  renderPositionCategorySummary(positions);
  renderPositionMobileCards(filtered, totalValue);

  if (!filtered.length) {
    elements.positionsBody.innerHTML = `<tr><td colspan="8" class="empty-state">暂无符合条件的标的</td></tr>`;
    return;
  }

  elements.positionsBody.innerHTML = filtered
    .map(({ stock, strategy }) => {
      const pnl = decisionPnlCell(stock);
      const health = decisionHealthCell(stock, totalValue);
      const sourceClass = marketKind(stock) === "HK" ? "hk" : "";
      const marginTone = decisionMarginTone(strategy.margin);
      const returnTone = strategy.shield.passed ? "core" : "reduce";
      const ownerTone = decisionToneClass(strategy.ownerAudit.tone);
      const category = positionCategory(stock, strategy);
      const stockMeta = decisionStockMeta(stock);
      const stockTooltip = `${stock.name} · ${stockMeta}`;

      return `
        <tr>
          <td>
            <div class="stock-cell">
              <span class="ticker ${sourceClass}">${escapeHTML(stock.symbol.slice(0, 4))}</span>
              <a class="stock-name stock-link" href="${stockHash(stock.symbol)}" title="${escapeHTML(stockTooltip)}">
                <strong>${escapeHTML(stock.name)}</strong>
                <span>${escapeHTML(stockMeta)}</span>
              </a>
            </div>
          </td>
          <td data-label="分类">
            ${positionCategoryPill(category)}
          </td>
          <td data-label="市值/现价">
            ${decisionMarketCell(stock)}
          </td>
          <td data-label="盈亏" class="${pnl.className}">
            ${pnl.html}
          </td>
          <td data-label="综合回报率">
            <span class="health-pill ${returnTone}">${displayDividendRatio(strategy.shield.value)} / ${displayDividendRatio(strategy.shield.target)}</span>
          </td>
          <td data-label="安全边际">
            <span class="health-pill ${marginTone}">${Number.isFinite(strategy.margin) ? percent(strategy.margin * 100, false) : "-"}</span>
          </td>
          <td data-label="长期评分">
            <span class="health-pill ${ownerTone}">${strategy.ownerAudit.hasAudit ? `${strategy.ownerAudit.score}/100` : "待评分"}</span>
          </td>
          <td data-label="健康状态">
            <div class="health-status-cell" title="${escapeHTML(health.title)}">
              <strong class="health-status-score ${health.tone}">${escapeHTML(health.scoreText)}</strong>
              <small class="health-status-weight">${escapeHTML(health.weightText)}</small>
            </div>
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
            <span style="width: ${holdingsMasked ? 0 : share}%; background: ${palette[index % palette.length]}"></span>
          </div>
          <span>${escapeHTML(privateText(`${share.toFixed(1)}%`))}</span>
        </a>
      `;
    })
    .join("");
}

function renderTrades() {
  const recentTrades = [...state.trades]
    .filter((trade) => normalizeAssetType(trade.assetType) !== "fund")
    .reverse();

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
          <strong>${escapeHTML(trade.symbol)} · ${sideText}</strong>
          <span class="${privateClass(sideClass)}">${escapeHTML(privateText(currency(trade.price, trade.currency)))}</span>
          <small>${trade.date} · ${escapeHTML(trade.name)}</small>
          <small>${escapeHTML(privateText(`${trade.shares} 股`))} · 最新价 ${escapeHTML(privateText(currency(trade.currentPrice, trade.currency)))}</small>
        </div>
      `;
    })
    .join("");
}

function parseValuationRangeText(text) {
  const numbers = String(text ?? "").match(/\d+(?:\.\d+)?/g)?.map(Number).filter(Number.isFinite) ?? [];
  if (!numbers.length) return null;
  if (numbers.length === 1) return { low: numbers[0], base: numbers[0], high: numbers[0] };
  const low = Math.min(numbers[0], numbers[1]);
  const high = Math.max(numbers[0], numbers[1]);
  return { low, base: (low + high) / 2, high };
}

function valuationRangeView(stock) {
  const explicitRange = stock?.valuation?.range;
  const currencyCode = String(explicitRange?.currency || stock?.valuation?.currency || stock?.currency || "CNY").toUpperCase();
  if (explicitRange && Number.isFinite(finiteNumber(explicitRange.low)) && Number.isFinite(finiteNumber(explicitRange.high))) {
    const base = finiteNumber(explicitRange.base) ?? ((explicitRange.low + explicitRange.high) / 2);
    const margin = finiteNumber(explicitRange.marginOfSafety) ?? (base > 0 && stock.currentPrice > 0 ? (base - stock.currentPrice) / base : null);
    return { low: explicitRange.low, base, high: explicitRange.high, margin, currency: currencyCode, source: "三情景假设" };
  }
  const parsed = parseValuationRangeText(stock?.fairValueRange);
  if (parsed) {
    const margin = parsed.base > 0 && stock.currentPrice > 0 ? (parsed.base - stock.currentPrice) / parsed.base : marginValue(stock);
    return { ...parsed, margin, currency: currencyCode, source: "公允区间文本" };
  }
  const intrinsicValue = finiteNumber(stock?.intrinsicValue);
  if (Number.isFinite(intrinsicValue) && intrinsicValue > 0) {
    const low = intrinsicValue * 0.85;
    const high = intrinsicValue * 1.15;
    return { low, base: intrinsicValue, high, margin: marginValue(stock), currency: currencyCode, source: "内在价值推导" };
  }
  return { low: null, base: null, high: null, margin: null, currency: currencyCode, source: "待补充" };
}

function fallbackValuationScenarios(stock) {
  if (Array.isArray(stock?.valuation?.scenarios) && stock.valuation.scenarios.length) return stock.valuation.scenarios;
  const latest = latestAnnualFinancial(stock);
  const valuation = financialValuation(stock);
  const revenueGrowth = finiteNumber(latest.revenueYoY) ?? compoundGrowth(stock, "revenue") ?? 0;
  const profitMargin = finiteNumber(latest.netMargin) ?? finiteNumber(latest.profitMargin) ?? 0;
  const fcf = finiteNumber(latest.freeCashFlow);
  const pe = finiteNumber(valuation.pe);
  return [
    { name: "保守", revenueGrowth: revenueGrowth * 0.5, profitMargin: profitMargin * 0.9, fcf: Number.isFinite(fcf) ? fcf * 0.85 : null, discountRate: 0.105, reasonablePe: Number.isFinite(pe) ? pe * 0.85 : null },
    { name: "基准", revenueGrowth, profitMargin, fcf, discountRate: 0.095, reasonablePe: Number.isFinite(pe) ? pe : null },
    { name: "乐观", revenueGrowth: revenueGrowth * 1.2, profitMargin: profitMargin * 1.05, fcf: Number.isFinite(fcf) ? fcf * 1.15 : null, discountRate: 0.085, reasonablePe: Number.isFinite(pe) ? pe * 1.1 : null }
  ];
}

function valuationScenarioCell(stock) {
  const currencyCode = stock?.valuation?.currency || latestAnnualFinancial(stock).currency || stock?.currency || "CNY";
  return `
    <div class="valuation-scenario-grid">
      ${fallbackValuationScenarios(stock).map((scenario) => `
        <div>
          <strong>${escapeHTML(scenario.name || "情景")}</strong>
          <span>收入 ${financialRatio(scenario.revenueGrowth)}</span>
          <span>利润率 ${financialRatio(scenario.profitMargin)}</span>
          <span>FCF ${Number.isFinite(finiteNumber(scenario.fcf)) ? financialAmount(scenario.fcf, currencyCode) : "-"}</span>
          <span>折现 ${financialRatio(scenario.discountRate)}</span>
          <span>PE ${financialMultiple(scenario.reasonablePe)} / FCF ${financialMultiple(scenario.reasonablePfcf)}</span>
        </div>
      `).join("")}
    </div>
  `;
}

function renderValuationModule(positions) {
  if (!elements.valuationModuleList) return;
  const stocks = auditUniverse(positions)
    .map((stock) => ({ ...stock, screening: screeningProfile(stock), valuationRange: valuationRangeView(stock) }))
    .sort((a, b) => compareNullableNumbers(a.valuationRange.margin, b.valuationRange.margin, "desc") || compareNullableStrings(a.name, b.name, "asc"));

  if (!stocks.length) {
    elements.valuationModuleList.innerHTML = `<div class="empty-state">暂无股票池标的</div>`;
    return;
  }

  elements.valuationModuleList.innerHTML = stocks.map((stock) => {
    const range = stock.valuationRange;
    const marginTone = decisionMarginTone(range.margin);
    const priceText = Number.isFinite(finiteNumber(stock.currentPrice)) ? currency(stock.currentPrice, stock.currency) : "-";
    const rangeText = Number.isFinite(finiteNumber(range.low)) && Number.isFinite(finiteNumber(range.high))
      ? `${currency(range.low, range.currency)} - ${currency(range.high, range.currency)}`
      : "估值区间待补充";
    return `
      <article class="valuation-card">
        <div class="valuation-card-head">
          <a class="stock-name stock-link" href="${stockHash(stock.symbol)}">
            <strong>${escapeHTML(stock.name)}</strong>
            <span>${escapeHTML(stock.symbol)} · ${escapeHTML(displayText(firstIndustry(stock.industry), "未分类"))}</span>
          </a>
          <span class="health-pill ${marginTone}">${Number.isFinite(range.margin) ? percent(range.margin * 100, false) : "安全边际待补"}</span>
        </div>
        <div class="valuation-range-grid">
          <div><span>现价</span><strong>${escapeHTML(privateText(priceText))}</strong></div>
          <div><span>估值区间</span><strong>${escapeHTML(privateText(rangeText))}</strong><small>${escapeHTML(range.source)}</small></div>
          <div><span>基准值</span><strong>${Number.isFinite(finiteNumber(range.base)) ? escapeHTML(privateText(currency(range.base, range.currency))) : "-"}</strong></div>
          <div><span>选股分</span><strong>${escapeHTML(String(stock.screening.score))}/100</strong><small>${escapeHTML(stock.screening.pass ? "通过硬否决" : stock.screening.rejects.join(" / "))}</small></div>
        </div>
        ${valuationScenarioCell(stock)}
        <div class="valuation-card-actions">
          <button class="ghost-button compact-link" type="button" data-create-hold-log="${escapeHTML(stock.symbol)}">记录继续持有</button>
        </div>
      </article>
    `;
  }).join("");
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

  if (elements.candidateList) {
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
      : `<div class="empty-state compact-empty">当前筛选下暂无跟踪标的</div>`;
  }
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

function isOverviewActionSignal(signal) {
  if (signal.stock.sourceType === "holding") return true;
  return signal.strategy.bucket === "main" || signal.strategy.bucket === "cigar";
}

function buildActionConclusion(positions) {
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const totalAssets = totalValue + (finiteNumber(state.cash) ?? 0);
  const cashRatio = totalAssets ? (finiteNumber(state.cash) ?? 0) / totalAssets : 0;
  const signals = buildOpportunitySignals(positions);
  const overviewSignals = signals.filter(isOverviewActionSignal);
  const buySignals = overviewSignals.filter(({ strategy }) => strategy.bucket === "main" || strategy.bucket === "cigar");
  const reduceSignals = overviewSignals.filter(({ stock, strategy }) => stock.sourceType === "holding" && strategy.bucket === "excluded");
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
  if (overviewSignals.length) {
    return {
      tone: "watch",
      status: "过渡观察",
      detail: "旧仓不强制卖出，新资金只等综合回报/安全边际或烟蒂条件达标",
      reasons: [overviewSignals[0].reasons[0], ...reasons].slice(0, 3)
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
  elements.actionConclusion.className = `overview-conclusion ${conclusion.tone}`;
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
    <a class="decision-action-row ${item.tone}" href="${stockHash(item.symbol)}">
      <div class="decision-action-stock">
        <span class="decision-action-rank">${index + 1}</span>
        <div>
          <strong>${escapeHTML(item.name)}</strong>
          <small>${escapeHTML(item.symbol)} · ${escapeHTML(item.type)}</small>
        </div>
      </div>
      <p>${escapeHTML(item.detail)}</p>
    </a>
  `;
}

function overviewCandidateMarker(stock, strategy) {
  const margin = finiteNumber(strategy?.margin) ?? calculatedMarginOfSafety(stock);
  if (!Number.isFinite(margin)) return 66;
  if (margin >= MAIN_DCF_MARGIN_TARGET) return 36;
  if (margin <= 0) return 76;
  return Math.max(36, Math.min(76, 66 - (margin / MAIN_DCF_MARGIN_TARGET) * 30));
}

function renderCockpitSignal(title, detail, marker = null) {
  return `
    <article class="cockpit-signal">
      <strong>${escapeHTML(title)}</strong>
      <span>${escapeHTML(detail)}</span>
      ${Number.isFinite(marker) ? `<div class="mini-range"><i style="left: ${marker.toFixed(0)}%"></i></div>` : ""}
    </article>
  `;
}

function renderOverviewBuyCandidates(positions) {
  if (!elements.overviewBuyCandidates) return;
  const candidates = decisionUniverse(positions)
    .filter((stock) => stock.sourceType !== "holding")
    .map((stock) => ({
      stock,
      screening: screeningProfile(stock),
      strategy: strategyProfile(stock),
      valuation: valuationRangeView(stock)
    }))
    .filter((item) => item.screening.pass && item.strategy.bucket !== "excluded")
    .sort((a, b) => {
      const marginA = finiteNumber(a.strategy.margin) ?? -Infinity;
      const marginB = finiteNumber(b.strategy.margin) ?? -Infinity;
      return b.screening.score - a.screening.score || marginB - marginA || a.stock.name.localeCompare(b.stock.name, "zh-CN");
    })
    .slice(0, 3);

  elements.overviewBuyCandidates.innerHTML = candidates.length
    ? candidates.map(({ stock, screening, strategy }) => {
        const margin = Number.isFinite(strategy.margin) ? `安全边际 ${percent(strategy.margin * 100, false)}` : "安全边际待补";
        const detail = `选股分 ${screening.score} · ${strategy.shield.source || "综合回报"} ${displayDividendRatio(strategy.shield.value)} · ${margin}`;
        return renderCockpitSignal(stock.name, detail, overviewCandidateMarker(stock, strategy));
      }).join("")
    : `<div class="empty-state compact-empty">暂无通过硬否决且接近买点的候选</div>`;
}

function renderOverviewRiskReview(positions) {
  if (!elements.overviewRiskReview) return;
  const totalAssets = positions.reduce((sum, item) => sum + (finiteNumber(item.marketValueCny) ?? 0), 0) + (finiteNumber(state.cash) ?? 0);
  const maxPosition = positions.reduce((max, item) => (item.marketValueCny > (max?.marketValueCny ?? 0) ? item : max), null);
  const riskSignal = buildOpportunitySignals(positions)
    .find((signal) => signal.stock.sourceType === "holding" && signal.strategy.bucket === "excluded");
  const qualityIssues = buildDataQualityIssues(positions);
  const warningCount = qualityIssues.filter((issue) => issue.tone === "warn" || issue.tone === "error").length;
  const logCount = Array.isArray(state.decisionLogs) ? state.decisionLogs.length : 0;
  const items = [
    {
      title: "单票风险",
      detail: maxPosition
        ? `${maxPosition.name} ${percent(totalAssets ? (maxPosition.marketValueCny / totalAssets) * 100 : 0, false)}，继续看行业集中和现金缓冲。`
        : "暂无持仓。"
    },
    {
      title: "买入逻辑偏离",
      detail: riskSignal
        ? `${riskSignal.stock.name}：${stockActionReasonText(riskSignal.stock, riskSignal.strategy).replace(/^卖出理由[:：]\s*/, "")}`
        : "暂无持仓触发重大风险排除。"
    },
    {
      title: "日志缺口",
      detail: warningCount
        ? `数据体检 ${warningCount} 项提醒；已有 ${logCount} 条决策日志，复盘时优先补理由。`
        : `已有 ${logCount} 条决策日志；新交易继续强制填写理由。`
    }
  ];

  elements.overviewRiskReview.innerHTML = items
    .map((item) => renderCockpitSignal(item.title, item.detail))
    .join("");
}

function renderCommitteeOverview(positions) {
  const actions = buildOpportunitySignals(positions)
    .filter(isOverviewActionSignal)
    .slice(0, 6)
    .map((signal) => ({
      tone: signal.tone,
      type: strategyBucketLabel(signal.strategy.bucket),
      symbol: signal.stock.symbol,
      name: signal.stock.name,
      meta: `${displayDividendRatio(signal.strategy.shield.value)} 综合回报 · 安全边际 ${Number.isFinite(signal.strategy.margin) ? percent(signal.strategy.margin * 100, false) : "-"} · ${signal.strategy.ownerAudit.text}`,
      detail: stockActionReasonText(signal.stock, signal.strategy)
    }));
  if (elements.decisionQueueCount) {
    elements.decisionQueueCount.textContent = actions.length ? `${actions.length} 项待判断` : "无待办";
  }
  elements.committeeConsensus.innerHTML = actions.length
    ? actions.map(renderExecutiveActionItem).join("")
    : `<div class="empty-state compact-empty">暂无需要立即处理的动作</div>`;
  renderOverviewBuyCandidates(positions);
  renderOverviewRiskReview(positions);
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
  const sourceText = stock.sourceType === "holding" ? "持仓" : "跟踪";
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
  const sourceText = stock.sourceType === "holding" ? "持仓" : "跟踪";
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

function masterMatrixFilterCount(items, key) {
  if (key === "holding") return items.filter((item) => item.stock.sourceType === "holding").length;
  if (key === "candidate") return items.filter((item) => item.stock.sourceType === "candidate").length;
  return items.length;
}

function renderMasterMatrixFilters(items) {
  if (!elements.masterMatrixFilters) return;
  elements.masterMatrixFilters.innerHTML = MASTER_MATRIX_FILTERS.map((filter) => {
    const active = masterMatrixFilter === filter.key;
    return `
      <button class="${active ? "active" : ""}" type="button" data-master-matrix-filter="${filter.key}" aria-pressed="${active ? "true" : "false"}">
        <span>${escapeHTML(filter.label)}</span>
        <strong>${masterMatrixFilterCount(items, filter.key)}</strong>
      </button>
    `;
  }).join("");
}

function filterMasterMatrixRows(items) {
  if (masterMatrixFilter === "holding") return items.filter((item) => item.stock.sourceType === "holding");
  if (masterMatrixFilter === "candidate") return items.filter((item) => item.stock.sourceType === "candidate");
  return items;
}

function renderMasterMatrix(positions) {
  if (!elements.masterMatrix) return;
  const universe = strategyUniverseItems(positions);
  renderMasterMatrixFilters(universe);
  const rows = filterMasterMatrixRows(universe)
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
        <span>净现比</span>
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
    : `<div class="empty-state compact-empty">当前索引下暂无可对照标的</div>`;
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
  const ratio = netCashMarketCapRatio(netCash);
  if (!Number.isFinite(ratio)) return "watch";
  if (ratio <= 0) return "risk";
  return ratio >= 0.2 ? "strong" : "watch";
}

function netCashApplicable(stock) {
  return !/(^|\/)(银行|保险|券商|证券|信托|财富管理)(\/|$)/.test(String(stock?.industry ?? ""));
}

function netCashMarketCapRatio(netCash) {
  const netCashCny = finiteNumber(netCash?.netCashCny);
  const marketCapCny = finiteNumber(netCash?.marketCapCny);
  if (!Number.isFinite(netCashCny) || !Number.isFinite(marketCapCny) || marketCapCny <= 0) return null;
  return netCashCny / marketCapCny;
}

function netCashMatrixText(stock, netCash) {
  if (!netCashApplicable(stock)) return "不适用";
  const ratio = netCashMarketCapRatio(netCash);
  if (!Number.isFinite(finiteNumber(netCash?.netCashCny))) return "未录入";
  if (!Number.isFinite(finiteNumber(netCash?.marketCapCny))) return "缺市值";
  return Number.isFinite(ratio) ? percent(ratio * 100, false) : "-";
}

function fcfMatrixText(stock, netCash) {
  if (!netCashApplicable(stock)) return "不适用";
  return financialMultiple(netCash?.fcfMultiple);
}

function firstIndustry(industry) {
  return String(industry ?? "").split("/").map((item) => item.trim()).find(Boolean) || "未分类";
}

function normalizeIndustryId(value) {
  return String(value ?? "")
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9\u4e00-\u9fa5]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

function industryRecords() {
  return Array.isArray(state.industries) ? state.industries : [];
}

function findIndustryRecord(id) {
  const normalized = normalizeIndustryId(id);
  return industryRecords().find((item) => normalizeIndustryId(item.id || item.name) === normalized) ?? null;
}

function industryMetricText(metric) {
  if (metric?.valueText) return metric.valueText;
  const value = finiteNumber(metric?.latestValue);
  if (Number.isFinite(value)) {
    return `${value}${metric?.unit ? ` ${metric.unit}` : ""}`;
  }
  return displayText(metric?.valueText, "待录入");
}

function industryMetricTone(metric) {
  const value = String(metric?.tone ?? "").trim().toLowerCase();
  if (/strong|positive|good|improve|up|改善|向好|增长/.test(value)) return "strong";
  if (/risk|negative|bad|down|下降|承压|恶化|风险/.test(value)) return "risk";
  if (/watch|mixed|neutral|观察|分化|待验证/.test(value)) return "watch";
  return "watch";
}

function metricSeries(metric) {
  return Array.isArray(metric?.series)
    ? metric.series
      .map((point) => ({
        date: String(point?.date ?? "").trim(),
        value: finiteNumber(point?.value)
      }))
      .filter((point) => point.date && Number.isFinite(point.value))
      .sort((a, b) => a.date.localeCompare(b.date))
    : [];
}

function chartValue(value) {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return "-";
  const abs = Math.abs(number);
  if (abs >= 1000) return number.toFixed(0);
  if (abs >= 100) return number.toFixed(1);
  if (abs >= 10) return number.toFixed(1);
  return number.toFixed(2);
}

function renderLineChart(series, options = {}) {
  const points = series.filter((point) => Number.isFinite(point.value));
  if (points.length < 2) {
    return `
      <div class="line-chart-empty">
        ${escapeHTML(points.length === 1 ? `近24个月只有 1 个有效点：${points[0].date} ${chartValue(points[0].value)}` : "近24个月趋势数据不足，等待更新。")}
      </div>
    `;
  }
  const width = options.width ?? 320;
  const height = options.height ?? 138;
  const left = options.left ?? 34;
  const right = options.right ?? 10;
  const top = options.top ?? 12;
  const bottom = options.bottom ?? 24;
  const values = points.map((point) => point.value);
  const min = Math.min(...values);
  const max = Math.max(...values);
  const span = max - min || Math.max(1, Math.abs(max) * 0.1);
  const yMin = min - span * 0.12;
  const yMax = max + span * 0.12;
  const plotWidth = width - left - right;
  const plotHeight = height - top - bottom;
  const coords = points.map((point, index) => {
    const x = left + (plotWidth * index) / (points.length - 1);
    const y = top + ((yMax - point.value) / (yMax - yMin)) * plotHeight;
    return { ...point, x, y };
  });
  const line = coords.map((point, index) => `${index === 0 ? "M" : "L"} ${point.x.toFixed(1)} ${point.y.toFixed(1)}`).join(" ");
  const area = `${line} L ${coords[coords.length - 1].x.toFixed(1)} ${top + plotHeight} L ${coords[0].x.toFixed(1)} ${top + plotHeight} Z`;
  const first = coords[0];
  const last = coords[coords.length - 1];
  const midY = top + plotHeight / 2;
  return `
    <svg viewBox="0 0 ${width} ${height}" role="img" aria-label="${escapeHTML(options.label || "趋势折线图")}">
      <line class="grid-line" x1="${left}" x2="${width - right}" y1="${top}" y2="${top}" />
      <line class="grid-line" x1="${left}" x2="${width - right}" y1="${midY}" y2="${midY}" />
      <line class="grid-line" x1="${left}" x2="${width - right}" y1="${top + plotHeight}" y2="${top + plotHeight}" />
      ${options.hideAxis ? "" : `
        <text class="chart-label" x="0" y="${top + 4}">${escapeHTML(chartValue(yMax))}</text>
        <text class="chart-label" x="0" y="${top + plotHeight + 3}">${escapeHTML(chartValue(yMin))}</text>
      `}
      <path class="chart-area" d="${area}" />
      <path class="chart-line" d="${line}" />
      <circle class="chart-dot" cx="${first.x.toFixed(1)}" cy="${first.y.toFixed(1)}" r="${options.dotRadius ?? 4}" />
      <circle class="chart-dot" cx="${last.x.toFixed(1)}" cy="${last.y.toFixed(1)}" r="${options.lastDotRadius ?? 4.5}" />
      ${options.hideLabels ? "" : `
        <text class="chart-value" x="${Math.max(left, last.x - 42).toFixed(1)}" y="${Math.max(14, last.y - 9).toFixed(1)}">${escapeHTML(chartValue(last.value))}</text>
        <text class="chart-label" x="${left}" y="${height - 4}">${escapeHTML(first.date)}</text>
        <text class="chart-label" text-anchor="end" x="${width - right}" y="${height - 4}">${escapeHTML(last.date)}</text>
      `}
    </svg>
  `;
}

function stockMatchesIndustry(stock, industry) {
  const explicitSymbols = new Set((industry.linkedSymbols ?? []).map(normalizeSymbol).filter(Boolean));
  if (explicitSymbols.has(normalizeSymbol(stock.symbol))) return true;

  const keywords = [industry.name, ...(industry.keywords ?? [])]
    .map((item) => String(item ?? "").trim())
    .filter(Boolean);
  if (!keywords.length) return false;

  const haystack = [stock.industry, stock.name, stock.symbol].filter(Boolean).join(" ");
  return keywords.some((keyword) => haystack.includes(keyword));
}

function industryStocks(industry, positions) {
  return decisionUniverse(positions)
    .filter((stock) => stockMatchesIndustry(stock, industry))
    .sort((a, b) => {
      const valueA = finiteNumber(a.marketValueCny) ?? 0;
      const valueB = finiteNumber(b.marketValueCny) ?? 0;
      return valueB - valueA || a.name.localeCompare(b.name, "zh-CN");
    });
}

function industryExposure(stocks) {
  return stocks
    .filter((stock) => stock.sourceType === "holding")
    .reduce((sum, stock) => sum + (finiteNumber(stock.marketValueCny) ?? 0), 0);
}

function industryCompanyAnalyses(industry, stocks) {
  const analyses = Array.isArray(industry?.companyAnalyses) ? industry.companyAnalyses : [];
  if (analyses.length) return analyses;
  return stocks.map((stock) => ({
    symbol: stock.symbol,
    name: stock.name,
    stance: stock.status,
    summary: stock.action || stock.notes || "由当前持仓/候选标的自动生成财报趋势卡。"
  }));
}

function renderIndustryMetric(metric) {
  const series = metricSeries(metric);
  const latest = series[series.length - 1] ?? null;
  return `
    <article class="industry-line-metric ${industryMetricTone(metric)}">
      <span>${escapeHTML(displayText(metric?.name, "未命名指标"))}</span>
      <strong>${escapeHTML(industryMetricText(metric))}</strong>
      <div class="industry-line-chart">
        ${renderLineChart(series, { label: displayText(metric?.name, "行业趋势折线图") })}
        <div class="chart-caption">
          <span>${escapeHTML(series.length >= 2 ? `近24个月 · ${series.length} 个有效点` : "近24个月 · 等待更多数据")}</span>
          <span>${escapeHTML(latest ? `最新 ${chartValue(latest.value)}${metric?.unit || ""}` : "最新值待更新")}</span>
        </div>
      </div>
      ${metric?.trendText ? `<p>${escapeHTML(metric.trendText)}</p>` : metric?.comment ? `<p>${escapeHTML(metric.comment)}</p>` : ""}
      <small>${escapeHTML([metric?.asOf, metric?.source].filter(Boolean).join(" · ") || "来源待补充")}</small>
    </article>
  `;
}

function renderIndustryStockChip(stock) {
  const sourceText = stock.sourceType === "holding" ? "持仓" : "跟踪";
  return `
    <a class="industry-stock-chip" href="${stockHash(stock.symbol)}">
      <strong>${escapeHTML(stock.name)}</strong>
      <span>${escapeHTML(stock.symbol)} · ${escapeHTML(sourceText)}</span>
    </a>
  `;
}

function industryTrendTone(direction) {
  const value = String(direction ?? "").trim().toLowerCase();
  if (!value) return "neutral";
  if (/strong|positive|improve|改善|扩张|增长|净现金|稳定/.test(value)) return "strong";
  if (/risk|negative|down|下降|承压|恶化|风险/.test(value)) return "risk";
  if (/watch|review|mixed|分化|观察|放缓|待验证/.test(value)) return "watch";
  return "neutral";
}

function annualFinancialSeries(stock, key) {
  return financialAnnuals(stock)
    .map((item) => ({
      date: String(item?.fiscalYear || item?.reportDate || "").slice(0, 4),
      value: finiteNumber(item?.[key])
    }))
    .filter((point) => point.date && Number.isFinite(point.value))
    .sort((a, b) => a.date.localeCompare(b.date));
}

function latestAndPreviousAnnual(stock) {
  const annual = financialAnnuals(stock);
  return {
    latest: annual[0] ?? {},
    previous: annual[1] ?? {},
    oldest: annual[annual.length - 1] ?? {}
  };
}

function yoyText(value, label = "同比") {
  const number = finiteNumber(value);
  return Number.isFinite(number) ? `${label}${number >= 0 ? "+" : ""}${financialRatio(number)}` : "";
}

function directionFromChange(current, previous, strong = 0.03, risk = -0.03) {
  const now = finiteNumber(current);
  const before = finiteNumber(previous);
  if (!Number.isFinite(now) || !Number.isFinite(before) || before === 0) return "watch";
  const change = (now - before) / Math.abs(before);
  if (change >= strong) return "strong";
  if (change <= risk) return "risk";
  return "watch";
}

function directionFromValue(value, strong, risk) {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return "watch";
  if (number >= strong) return "strong";
  if (number <= risk) return "risk";
  return "watch";
}

function computedCompanyTrends(stock) {
  if (!stock || !financialAnnuals(stock).length) return [];
  const { latest, previous, oldest } = latestAndPreviousAnnual(stock);
  const currencyCode = latestFinancialCurrency(stock);
  const dividendYield = calculatedDividendYield(stock);
  const margin = calculatedMarginOfSafety(stock) ?? finiteNumber(stock.marginOfSafety);
  const trends = [];

  const revenueSeries = annualFinancialSeries(stock, "revenue");
  trends.push({
    label: "收入",
    value: `${financialAmount(oldest.revenue, currencyCode)} → ${financialAmount(latest.revenue, currencyCode)}`,
    direction: Number.isFinite(finiteNumber(latest.revenueYoY)) ? directionFromValue(latest.revenueYoY, 0.05, 0) : directionFromChange(latest.revenue, previous.revenue),
    note: yoyText(latest.revenueYoY) || "多年财报口径自动生成。",
    series: revenueSeries
  });

  trends.push({
    label: "归母利润",
    value: `${financialAmount(oldest.netProfit, currencyCode)} → ${financialAmount(latest.netProfit, currencyCode)}`,
    direction: Number.isFinite(finiteNumber(latest.netProfitYoY)) ? directionFromValue(latest.netProfitYoY, 0.05, 0) : directionFromChange(latest.netProfit, previous.netProfit),
    note: yoyText(latest.netProfitYoY) || "对比最新年度和上年利润变化。",
    series: annualFinancialSeries(stock, "netProfit")
  });

  trends.push({
    label: "OCF / 收入",
    value: financialRatio(latest.operatingCashFlowToRevenue),
    direction: directionFromValue(latest.operatingCashFlowToRevenue, 0.15, 0.08),
    note: previous.operatingCashFlowToRevenue
      ? `上年 ${financialRatio(previous.operatingCashFlowToRevenue)}，观察现金转化。`
      : "观察收入是否真实转化为经营现金流。",
    series: annualFinancialSeries(stock, "operatingCashFlowToRevenue")
  });

  trends.push({
    label: "自由现金流",
    value: financialAmount(latest.freeCashFlow, currencyCode),
    direction: directionFromValue(latest.freeCashFlow, 0, -1),
    note: previous.freeCashFlow ? `上年 ${financialAmount(previous.freeCashFlow, currencyCode)}。` : "用于验证分红和防守属性。",
    series: annualFinancialSeries(stock, "freeCashFlow")
  });

  trends.push({
    label: "毛利率",
    value: financialRatio(latest.grossMargin),
    direction: directionFromChange(latest.grossMargin, previous.grossMargin, 0.02, -0.02),
    note: previous.grossMargin ? `上年 ${financialRatio(previous.grossMargin)}，观察价格、成本和结构。` : "观察价格、成本和产品结构变化。",
    series: annualFinancialSeries(stock, "grossMargin")
  });

  trends.push({
    label: "ROE / ROIC",
    value: `${financialRatio(latest.roe)} / ${financialRatio(latest.roic)}`,
    direction: directionFromValue(Math.min(finiteNumber(latest.roe) ?? 0, finiteNumber(latest.roic) ?? 0), 0.12, 0.06),
    note: "资本回报用于判断是否仍是高质量资产。",
    series: annualFinancialSeries(stock, "roic")
  });

  trends.push({
    label: "净现金",
    value: financialAmount(latest.netCash, currencyCode),
    direction: directionFromValue(latest.netCash, 0, -1),
    note: "狭义现金及短投减有息债务，需结合行业属性折价。",
    series: annualFinancialSeries(stock, "netCash")
  });

  trends.push({
    label: "股息率 / 安全边际",
    value: `${displayDividendRatio(dividendYield)} / ${Number.isFinite(margin) ? percent(margin * 100, false) : "-"}`,
    direction: Number.isFinite(margin) && margin >= 0.15 ? "strong" : Number.isFinite(margin) && margin < 0.05 ? "risk" : "watch",
    note: "回报和买点纪律由当前价格、分红和内在价值实时计算。",
    series: []
  });

  return trends;
}

function renderCompanySparkline(series) {
  const points = Array.isArray(series) ? series.filter((point) => Number.isFinite(point.value)) : [];
  if (points.length < 2) return "";
  return `
    <div class="industry-sparkline">
      ${renderLineChart(points, { width: 180, height: 46, left: 4, right: 4, top: 5, bottom: 5, hideAxis: true, hideLabels: true, dotRadius: 2.5, lastDotRadius: 3, label: "年度趋势" })}
    </div>
  `;
}

function renderIndustryCompanyAnalysis(analysis, stocks) {
  const symbol = normalizeSymbol(analysis?.symbol);
  const stock = stocks.find((item) => normalizeSymbol(item.symbol) === symbol);
  const name = displayText(analysis?.name, stock?.name || symbol || "未命名标的");
  const trends = computedCompanyTrends(stock);
  const judgments = Array.isArray(analysis?.judgments) ? analysis.judgments : [];
  const watchpoints = Array.isArray(analysis?.watchpoints) ? analysis.watchpoints : [];
  const sourceText = stock?.sourceType === "holding" ? "持仓" : stock?.sourceType === "candidate" ? "跟踪" : "未关联持仓";

  return `
    <article class="industry-company-card">
      <div class="industry-company-head">
        <div>
          <span>${escapeHTML([symbol, sourceText].filter(Boolean).join(" · "))}</span>
          ${symbol ? `<a href="${stockHash(symbol)}">${escapeHTML(name)}</a>` : `<strong>${escapeHTML(name)}</strong>`}
        </div>
        <em>${escapeHTML(displayText(analysis?.stance, stock?.status || "待判断"))}</em>
      </div>
      <p>${escapeHTML(displayText(analysis?.summary, "后续补充该标的趋势和定性判断。"))}</p>
      ${trends.length ? `
        <div class="industry-trend-grid">
          ${trends.map((trend) => `
            <div class="industry-trend ${industryTrendTone(trend?.direction)}">
              <span>${escapeHTML(displayText(trend?.label, "趋势"))}</span>
              <strong>${escapeHTML(displayText(trend?.value, "-"))}</strong>
              ${renderCompanySparkline(trend?.series)}
              ${trend?.note ? `<small>${escapeHTML(trend.note)}</small>` : ""}
            </div>
          `).join("")}
        </div>
      ` : `<div class="line-chart-empty">未找到多年财报数据，先保留定性判断。</div>`}
      ${judgments.length ? `
        <div class="industry-company-block">
          <span>定性判断</span>
          <ul>
            ${judgments.map((item) => `<li>${escapeHTML(item)}</li>`).join("")}
          </ul>
        </div>
      ` : ""}
      ${watchpoints.length ? `
        <div class="industry-company-block watchpoints">
          <span>后续观察</span>
          <ul>
            ${watchpoints.map((item) => `<li>${escapeHTML(item)}</li>`).join("")}
          </ul>
        </div>
      ` : ""}
    </article>
  `;
}

function researchDeskSectionVisible(section) {
  return researchDeskFilter === "all" || researchDeskFilter === section;
}

function researchDeskDisplayDate(value) {
  const text = String(value ?? "").trim();
  if (!text) return "";
  return text.slice(0, 10);
}

function researchDeskLatestDate(industries, updates) {
  return [
    ...industries.flatMap((industry) => [industry.updatedAt, industry.metricsUpdatedAt]),
    ...updates.map((item) => item.date)
  ]
    .map(researchDeskDisplayDate)
    .filter(Boolean)
    .sort()
    .at(-1) ?? "";
}

function buildResearchOpportunities() {
  return sortedCandidates().slice(0, 8).map((candidate) => {
    const strategy = strategyProfile(candidate);
    const distance = candidateBuyDistance(candidate);
    const confidence = confidenceMeta(candidate);
    const plan = findPlanForStock(candidate);
    return { candidate, strategy, distance, confidence, plan };
  });
}

function researchQueueFromDataIssues(issues) {
  return issues
    .filter((issue) => issue.tone === "warn" || issue.tone === "error")
    .slice(0, 5)
    .map((issue) => ({
      type: "data",
      tone: issue.tone === "error" ? "risk" : "watch",
      symbol: issue.symbol,
      name: issue.name || issue.symbol || "数据状态",
      label: "数据体检",
      title: issue.title,
      detail: issue.detail
    }));
}

function researchQueueFromOpportunities(opportunities) {
  return opportunities
    .filter(({ candidate, strategy, distance, confidence }) => {
      const margin = calculatedMarginOfSafety(candidate) ?? finiteNumber(candidate.marginOfSafety);
      return strategy.bucket === "main" ||
        strategy.shield.passed ||
        (Number.isFinite(margin) && margin >= MAIN_DCF_MARGIN_TARGET) ||
        (Number.isFinite(distance) && distance <= BUY_PROXIMITY) ||
        confidence.tone === "risk";
    })
    .slice(0, 5)
    .map(({ candidate, strategy, distance, confidence, plan }) => ({
      type: "opportunity",
      tone: strategy.bucket === "excluded" || confidence.tone === "risk" ? "risk" : strategy.tone,
      symbol: candidate.symbol,
      name: candidate.name,
      label: "候选机会",
      title: Number.isFinite(distance) && distance <= BUY_PROXIMITY ? "接近买点" : strategy.status,
      detail: strategy.blockers.length ? strategy.blockers.join("；") : displayText(plan?.advice, candidate.action)
    }));
}

function researchQueueFromFinancials(positions) {
  return auditUniverse(positions)
    .filter((stock) => !String(stock.financials?.updatedAt ?? "").trim())
    .slice(0, 4)
    .map((stock) => ({
      type: "financials",
      tone: "info",
      symbol: stock.symbol,
      name: stock.name,
      label: "财务更新",
      title: "缺财务更新时间",
      detail: "详情页可点击“更新财务”补齐多年财务和估值口径。"
    }));
}

function buildResearchQueue(positions, issues, opportunities) {
  const items = [
    ...researchQueueFromDataIssues(issues),
    ...researchQueueFromOpportunities(opportunities),
    ...researchQueueFromFinancials(positions)
  ];
  const seen = new Set();
  return items.filter((item) => {
    const key = `${item.type}:${normalizeSymbol(item.symbol) || item.name}:${item.title}`;
    if (seen.has(key)) return false;
    seen.add(key);
    return true;
  }).slice(0, 10);
}

function buildResearchUpdates(positions) {
  const stocks = auditUniverse(positions);
  const stockMap = new Map(stocks.map((stock) => [normalizeSymbol(stock.symbol), stock]));
  const updates = [];

  stocks.forEach((stock) => {
    (stock.researchUpdates ?? []).forEach((item) => {
      updates.push({
        symbol: stock.symbol,
        name: stock.name,
        date: item.importedAt || item.asOf || item.date,
        title: item.title || item.event?.title || "研究更新",
        detail: item.summary || item.notesAppend || item.event?.summary || displayText(stock.action, stock.status),
        source: "研究导入"
      });
    });
  });

  (state.decisionLogs ?? [])
    .filter((log) => log.type === "research")
    .forEach((log) => {
      const symbol = normalizeSymbol(log.symbol);
      const stock = stockMap.get(symbol);
      updates.push({
        symbol,
        name: log.name || stock?.name || symbol,
        date: log.date || log.createdAt,
        title: log.decision || "导入分析",
        detail: log.detail || log.discipline || "",
        source: "决策日志"
      });
    });

  return updates
    .filter((item) => item.date || item.title || item.detail)
    .sort((a, b) => String(b.date ?? "").localeCompare(String(a.date ?? "")))
    .slice(0, 8);
}

function researchDeskMetric(label, value, detail = "") {
  return `
    <article class="research-desk-metric">
      <span>${escapeHTML(label)}</span>
      <strong>${escapeHTML(value)}</strong>
      <small>${escapeHTML(detail)}</small>
    </article>
  `;
}

function researchDeskFilters(counts) {
  return `
    <div class="research-desk-filter" aria-label="研究台筛选">
      ${RESEARCH_DESK_FILTERS.map((filter) => {
        const active = filter.value === researchDeskFilter;
        const count = counts[filter.value] ?? 0;
        return `
          <button class="${active ? "active" : ""}" type="button" data-research-desk-filter="${filter.value}" aria-pressed="${active ? "true" : "false"}">
            <span>${escapeHTML(filter.label)}</span>
            <strong>${count}</strong>
          </button>
        `;
      }).join("")}
    </div>
  `;
}

function researchDeskSection(id, eyebrow, title, body, count = "") {
  const expanded = expandedResearchDeskSections.has(id);
  return `
    <section class="research-desk-section ${expanded ? "" : "is-collapsed"}" data-research-desk-section="${escapeHTML(id)}">
      <button class="research-desk-section-head" type="button" data-research-desk-toggle="${escapeHTML(id)}" aria-expanded="${expanded ? "true" : "false"}">
        <span>
          <small>${escapeHTML(eyebrow)}</small>
          <strong>${escapeHTML(title)}</strong>
        </span>
        <em>${escapeHTML(count ? String(count) : expanded ? "收起" : "展开")}</em>
      </button>
      <div class="research-desk-section-body">
        ${body}
      </div>
    </section>
  `;
}

function renderResearchQueueItem(item) {
  const content = `
    <div class="research-queue-type">${escapeHTML(item.label)}</div>
    <div class="research-queue-main">
      <strong>${escapeHTML(item.name || item.symbol || "未命名")}</strong>
      <span>${escapeHTML(item.title)}</span>
      <small>${escapeHTML(item.detail || "等待补充下一步动作。")}</small>
    </div>
  `;
  return item.symbol
    ? `<a class="research-queue-item ${item.tone}" href="${stockHash(item.symbol)}">${content}</a>`
    : `<div class="research-queue-item ${item.tone}">${content}</div>`;
}

function renderResearchOpportunityCard(item) {
  const { candidate, strategy, distance, confidence, plan } = item;
  const margin = calculatedMarginOfSafety(candidate) ?? finiteNumber(candidate.marginOfSafety);
  const score = strategy.ownerAudit.hasAudit ? `${strategy.ownerAudit.score}/100` : "-";
  return `
    <a class="research-opportunity-card" href="${stockHash(candidate.symbol)}">
      <div class="research-card-head">
        <div>
          <strong>${escapeHTML(candidate.name)}</strong>
          <span>${escapeHTML(candidate.symbol)} · ${escapeHTML(displayText(candidate.industry, "未分类"))}</span>
        </div>
        <em>${escapeHTML(strategy.status)}</em>
      </div>
      <div class="research-card-metrics">
        <span>综合回报 <strong>${escapeHTML(privateText(displayDividendRatio(strategy.shield.value)))}</strong></span>
        <span>安全边际 <strong>${escapeHTML(privateText(Number.isFinite(margin) ? percent(margin * 100, false) : "-"))}</strong></span>
        <span>买点距离 <strong>${escapeHTML(privateText(Number.isFinite(distance) ? displayBuyDistance(candidate) : "-"))}</strong></span>
        <span>长期评分 <strong>${escapeHTML(score)}</strong></span>
      </div>
      <p>${escapeHTML(strategy.blockers.length ? strategy.blockers.join("；") : displayText(plan?.advice, candidate.action))}</p>
      <small>${escapeHTML(confidence.text)} · ${escapeHTML(strategy.shield.source)}</small>
    </a>
  `;
}

function renderResearchDataIssue(issue) {
  const tag = issue.sourceType === "holding" ? "持仓" : issue.sourceType === "candidate" ? "跟踪" : issue.sourceType === "data" ? "数据" : "Plan";
  const content = `
    <div class="research-card-head">
      <div>
        <strong>${escapeHTML(issue.title)}</strong>
        <span>${escapeHTML(issue.name || issue.symbol || "数据状态")} ${issue.symbol ? `· ${escapeHTML(issue.symbol)}` : ""}</span>
      </div>
      <em>${escapeHTML(tag)}</em>
    </div>
    <p>${escapeHTML(issue.detail)}</p>
  `;
  return issue.symbol
    ? `<a class="research-data-issue ${issue.tone}" href="${stockHash(issue.symbol)}">${content}</a>`
    : `<div class="research-data-issue ${issue.tone}">${content}</div>`;
}

function renderResearchIndustryCard(industry, positions) {
  const id = normalizeIndustryId(industry.id || industry.name);
  const stocks = industryStocks(industry, positions);
  const exposure = industryExposure(stocks);
  const metrics = Array.isArray(industry.metrics) ? industry.metrics : [];
  const companyAnalyses = industryCompanyAnalyses(industry, stocks);
  const latestNote = Array.isArray(industry.notes) ? industry.notes[0] : null;
  return `
    <a class="research-industry-card" href="${industryHash(id)}">
      <div class="research-card-head">
        <div>
          <span>${escapeHTML(displayText(industry.category, "行业档案"))}</span>
          <strong>${escapeHTML(displayText(industry.name, id))}</strong>
        </div>
        <em>${escapeHTML(displayText(industry.status, "待完善"))}</em>
      </div>
      <p>${escapeHTML(displayText(industry.summary, "后续补充行业分析、关键指标和跟踪记录。"))}</p>
      <div class="research-card-metrics">
        <span>相关标的 <strong>${stocks.length}</strong></span>
        <span>持仓暴露 <strong>${escapeHTML(privateText(currency(exposure)))}</strong></span>
        <span>趋势指标 <strong>${metrics.length}</strong></span>
        <span>标的分析 <strong>${companyAnalyses.length}</strong></span>
      </div>
      <small>${escapeHTML(latestNote ? `${latestNote.date || "未记录日期"} · ${latestNote.title || latestNote.summary || "研究记录"}` : "暂无研究记录")}</small>
    </a>
  `;
}

function renderResearchUpdateCard(item) {
  const content = `
    <div class="research-card-head">
      <div>
        <strong>${escapeHTML(item.title || "研究更新")}</strong>
        <span>${escapeHTML(item.name || item.symbol || "未关联标的")}</span>
      </div>
      <em>${escapeHTML(researchDeskDisplayDate(item.date) || item.source)}</em>
    </div>
    <p>${escapeHTML(item.detail || "暂无摘要。")}</p>
    <small>${escapeHTML(item.source)}</small>
  `;
  return item.symbol
    ? `<a class="research-update-card" href="${stockHash(item.symbol)}">${content}</a>`
    : `<div class="research-update-card">${content}</div>`;
}

function renderIndustryDesk(positions) {
  if (!elements.industryList) return;
  const industries = industryRecords()
    .slice()
    .sort((a, b) => String(b.updatedAt ?? "").localeCompare(String(a.updatedAt ?? "")) || String(a.name ?? "").localeCompare(String(b.name ?? ""), "zh-CN"));
  const issues = buildDataQualityIssues(positions);
  const warningCount = issues.filter((issue) => issue.tone === "warn" || issue.tone === "error").length;
  const opportunities = buildResearchOpportunities();
  const queue = buildResearchQueue(positions, issues, opportunities);
  const updates = buildResearchUpdates(positions);
  const latestDate = researchDeskLatestDate(industries, updates);
  const counts = {
    queue: queue.length,
    opportunity: opportunities.length,
    data: issues.length,
    industry: industries.length
  };
  counts.all = counts.queue + counts.opportunity + counts.data + counts.industry;

  const queueSection = researchDeskSection(
    "queue",
    "Research Queue",
    "今日研究队列",
    queue.length
      ? `<div class="research-queue-list">${queue.map(renderResearchQueueItem).join("")}</div>`
      : `<div class="empty-state compact-empty">当前没有高优先级研究任务</div>`,
    queue.length
  );
  const opportunitySection = researchDeskSection(
    "opportunity",
    "Candidate Signals",
    "候选机会",
    opportunities.length
      ? `<div class="research-opportunity-grid">${opportunities.map(renderResearchOpportunityCard).join("")}</div>`
      : `<div class="empty-state compact-empty">当前暂无候选机会</div>`,
    opportunities.length
  );
  const dataSection = researchDeskSection(
    "data",
    "Data Quality",
    "数据体检",
    issues.length
      ? `<div class="research-data-list">${issues.slice(0, 12).map(renderResearchDataIssue).join("")}</div>`
      : `<div class="empty-state compact-empty">关键研究数据完整</div>`,
    warningCount ? `${warningCount} 项` : "通过"
  );
  const industrySection = researchDeskSection(
    "industry",
    "Industry Themes",
    "行业主题",
    industries.length
      ? `<div class="research-industry-grid">${industries.map((industry) => renderResearchIndustryCard(industry, positions)).join("")}</div>`
      : `<div class="empty-state compact-empty">暂无行业档案。可在 data/industries/ 下新增行业 JSON。</div>`,
    industries.length
  );
  const updatesSection = researchDeskSection(
    "updates",
    "Recent Updates",
    "最近研究更新",
    updates.length
      ? `<div class="research-update-list">${updates.map(renderResearchUpdateCard).join("")}</div>`
      : `<div class="empty-state compact-empty">暂无研究更新记录</div>`,
    updates.length
  );

  const dataInPrimary = researchDeskFilter === "data";
  const primarySections = [
    researchDeskSectionVisible("queue") ? queueSection : "",
    researchDeskSectionVisible("opportunity") ? opportunitySection : "",
    researchDeskSectionVisible("industry") ? industrySection : "",
    dataInPrimary ? dataSection : ""
  ].join("");
  const secondarySections = [
    !dataInPrimary && researchDeskSectionVisible("data") ? dataSection : "",
    researchDeskFilter === "all" ? updatesSection : ""
  ].join("");

  elements.industryList.innerHTML = `
    <div class="research-desk">
      <section class="research-desk-hero">
        <div class="research-desk-hero-main">
          <p class="eyebrow">Research Command</p>
          <h2>研究总控台</h2>
          <span>把行业主题、跟踪机会、数据体检和近期研究更新放在同一个入口，优先处理会影响买入纪律和复盘结论的事项。</span>
        </div>
        <div class="research-desk-actions">
          <button class="primary-button" type="button" data-open-research>
            <span>＋</span>
            导入分析
          </button>
          <button class="ghost-button" type="button" data-update-industries>
            <span>↻</span>
            更新行业数据
          </button>
        </div>
        <div class="research-desk-metrics">
          ${researchDeskMetric("行业档案", `${industries.length}`, "主题研究")}
          ${researchDeskMetric("跟踪标的", `${state.candidates?.length ?? 0}`, "晴仓30候选")}
          ${researchDeskMetric("待处理问题", warningCount ? `${warningCount} 项` : "通过", issues.length ? `总计 ${issues.length}` : "关键数据完整")}
          ${researchDeskMetric("最近更新", latestDate || "-", latestDate ? "研究/行业记录" : "暂无记录")}
        </div>
      </section>

      ${researchDeskFilters(counts)}

      <div class="research-desk-grid ${secondarySections ? "" : "single-column"}">
        <main class="research-desk-primary">
          ${primarySections || `<div class="empty-state compact-empty">当前筛选下暂无内容</div>`}
        </main>
        <aside class="research-desk-secondary">
          ${secondarySections}
        </aside>
      </div>
    </div>
  `;
}

function renderIndustryDetail(positions, id) {
  if (!elements.industryDetail) return;
  const industry = findIndustryRecord(id);
  if (!industry) {
    elements.industryDetail.innerHTML = `
      <section class="panel">
        <div class="empty-state">未找到该行业档案</div>
      </section>
    `;
    return;
  }

  const stocks = industryStocks(industry, positions);
  const keyQuestions = Array.isArray(industry.keyQuestions) ? industry.keyQuestions : [];
  const companyAnalyses = industryCompanyAnalyses(industry, stocks);
  const metrics = Array.isArray(industry.metrics) ? industry.metrics : [];
  const exposure = industryExposure(stocks);
  const updateText = [industry.metricsUpdatedAt, industry.updatedAt].map((item) => String(item ?? "").trim()).find(Boolean) || "";

  elements.industryDetail.innerHTML = `
    <section class="industry-detail-hero">
      <a class="ghost-button detail-back" href="#industry">返回研究台</a>
      <div>
        <p class="eyebrow">${escapeHTML(displayText(industry.category, "Industry"))}</p>
        <h2>${escapeHTML(displayText(industry.name, "未命名行业"))}</h2>
      </div>
      <div class="detail-hero-meta">
        <span>${escapeHTML(displayText(industry.status, "待完善"))}</span>
        <small>${escapeHTML(updateText ? `更新 ${updateText}` : "更新时间待补充")}</small>
      </div>
    </section>

    <section class="metrics-grid industry-metrics-summary">
      ${metricCard("持仓暴露", privateText(currency(exposure)), "按当前市值折人民币")}
      ${metricCard("标的分析", `${companyAnalyses.length}`, "趋势与定性判断")}
      ${metricCard("覆盖公司", stocks.map((stock) => stock.name).join(" / ") || "-", "行业内当前关注对象")}
      ${metricCard("趋势指标", `${metrics.length}`, industry.metricsUpdatedAt ? "runtime 行业数据" : "等待行业数据更新")}
    </section>

    <section class="panel industry-summary-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Thesis</p>
          <h2>行业结论</h2>
        </div>
        <button class="primary-button" type="button" data-update-industries>
          <span>↻</span>
          更新行业数据
        </button>
      </div>
      <div class="industry-summary-body">
        <p>${escapeHTML(displayText(industry.summary, "后续补充行业景气、供需、成本、政策、估值和代表公司分析。"))}</p>
        <div class="industry-question-list">
          ${keyQuestions.length
            ? keyQuestions.map((item) => `<span>${escapeHTML(item)}</span>`).join("")
            : `<span>后续补充关键研究问题</span>`}
        </div>
        ${industry.discipline ? `
          <div class="industry-discipline">
            <span>执行纪律</span>
            <p>${escapeHTML(industry.discipline)}</p>
          </div>
        ` : ""}
      </div>
    </section>

    <section class="panel industry-trend-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">External Metrics</p>
          <h2>行业趋势数据</h2>
        </div>
        <small>${escapeHTML(industry.metricsUpdatedAt ? `runtime 更新 ${industry.metricsUpdatedAt}` : "尚未生成 runtime 行业数据")}</small>
      </div>
      <div class="industry-line-metric-grid">
        ${metrics.length
          ? metrics.map(renderIndustryMetric).join("")
          : `<div class="empty-state compact-empty">暂无行业趋势数据。点击“更新行业数据”后写入 data/runtime/industry_metrics.json。</div>`}
      </div>
    </section>

    ${companyAnalyses.length ? `
      <section class="panel industry-company-panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Company Trends</p>
            <h2>标的趋势分析</h2>
          </div>
        </div>
        <div class="industry-company-grid">
          ${companyAnalyses.map((analysis) => renderIndustryCompanyAnalysis(analysis, stocks)).join("")}
        </div>
      </section>
    ` : ""}
  `;
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
        <strong>${escapeHTML(privateText(item.value))}</strong>
        <small>${escapeHTML(privateText(item.detail))}</small>
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
  const candidates = (state.candidates ?? [])
    .filter((candidate) => {
      const symbol = normalizeSymbol(candidate.symbol);
      return symbol && !holdingSymbols.has(symbol);
    })
    .map((candidate) => ({
      ...candidate,
      sourceType: "candidate"
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

  (state.dataStatus?.issues ?? []).forEach((issue) => {
    pushIssue(issue.tone || "info", { name: "数据持久化", sourceType: "data" }, issue.title || "数据状态提醒", issue.detail || "");
  });

  stocks.forEach((stock) => {
    const currentPrice = finiteNumber(stock.currentPrice);
    const previousClose = finiteNumber(stock.previousClose);
    const intrinsicValue = finiteNumber(stock.intrinsicValue);
    const qualityScore = finiteNumber(stock.qualityScore);
    const dividend = stock.dividend;
    const dividendYield = calculatedDividendYield(stock);
    const lagDays = dateDiffDays(stock.currentPriceDate, referenceDate);

    if (stock.duplicateHolding) {
      pushIssue("warn", stock, "跟踪标的重复持仓", "该标的已有持仓，晴仓30展示会去重；建议只保留一处维护。");
    }
    if (!Number.isFinite(currentPrice) || currentPrice <= 0) {
      pushIssue("warn", stock, "缺最新价", "会影响市值、安全边际、行动列表和综合回报率计算。");
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
        pushIssue("info", stock, "缺股息数据", "预计年股息不会计入该持仓。");
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
      const tag = issue.sourceType === "holding" ? "持仓" : issue.sourceType === "candidate" ? "跟踪" : issue.sourceType === "data" ? "数据" : "Plan";
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
  renderDisciplineDashboard(positions);
  renderDataQuality(positions);
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
  const taxNote = dividendTaxText(stock);
  const taxAdjusted = dividendTaxRate(stock) > 0;

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
          <div><span>每股分红</span><strong>${Number.isFinite(perShare) ? currency(perShare, currencyCode) : "-"}</strong><small>${taxAdjusted ? "税前公告口径" : ""}</small></div>
          <div><span>${taxAdjusted ? "税后股息率" : "股息率"}</span><strong>${displayDividendRatio(dividendYield)}</strong><small>${escapeHTML(taxNote)}</small></div>
          <div><span>回报门槛</span><strong>${displayDividendRatio(shield.target)}</strong><small>${marketKind(stock) === "HK" ? "H股主策略" : "A股主策略"}</small></div>
          <div><span>${taxAdjusted ? "税后综合回报率" : "综合回报率"}</span><strong>${displayDividendRatio(shareholderReturnYield)}</strong><small>${taxAdjusted ? "回购不折税" : ""}</small></div>
          <div><span>股息可靠性</span><strong>${badge(reliability.text, reliability.tone)}</strong></div>
          <div><span>预估财年</span><strong>${escapeHTML(displayText(dividend.forecastFiscalYear, "-"))}</strong></div>
          <div><span>预估每股</span><strong>${Number.isFinite(forecastPerShare) ? currency(forecastPerShare, forecastCurrency) : "-"}</strong></div>
          <div><span>${taxAdjusted ? "预估税后股息率" : "预估股息率"}</span><strong>${displayDividendRatio(forecastYield)}</strong><small>参考，不作为主策略门槛</small></div>
          <div><span>${taxAdjusted ? "预估年现金（税前）" : "预估年现金"}</span><strong>${Number.isFinite(annualLocal) ? currency(annualLocal, currencyCode) : Number.isFinite(estimatedCash) ? currency(estimatedCash, currencyCode) : "-"}</strong><small>${isHolding && annualCny > 0 ? `折人民币税后 ${currency(annualCny)}` : ""}</small></div>
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
  const taxAdjusted = dividendTaxRate(stock) > 0;
  const taxText = taxAdjusted ? dividendTaxText(stock) : displayText(dividend.taxNote, "未折扣分红税");
  const dividendFormula = hasCashDividendFormula
    ? `${taxAdjusted ? "税后" : ""}现金分红总额 / 总市值`
    : dividend.dividendPerShare
      ? `${taxAdjusted ? "税后" : ""}每股分红 / 最新价`
      : "暂无";
  const returnFormula = hasBuyback ? `(${taxAdjusted ? "税后" : ""}现金分红总额 + 回购金额) / 总市值` : "同股息率";

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
          <small>${escapeHTML(`${displayText(dividend.fiscalYear, "财年未知")} · ${taxText}`)}</small>
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

function renderFinancialsPanel(stock, options = {}) {
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

  const tableMarkup = options.collapsibleTable
    ? `
      <details class="financial-table-disclosure">
        <summary>多年财务表</summary>
        ${renderFinancialTable(stock, annual, currencyCode)}
      </details>
    `
    : renderFinancialTable(stock, annual, currencyCode);

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
        ${tableMarkup}
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

function stockDetailNav() {
  return `
    <nav class="detail-section-nav stock-detail-nav" aria-label="详情分段导航">
      ${STOCK_DETAIL_NAV_ITEMS.map((item) => `
        <button class="${activeStockDetailSection === item.id ? "active" : ""}" type="button" data-detail-section="${item.id}" aria-pressed="${activeStockDetailSection === item.id ? "true" : "false"}">
          <span class="desktop-label">${escapeHTML(item.desktopLabel)}</span>
          <span class="mobile-label">${escapeHTML(item.label)}</span>
        </button>
      `).join("")}
    </nav>
  `;
}

function stockDetailMetricButton(target, label, value, detail = "", className = "") {
  return `
    <button class="stock-detail-signal-card" type="button" data-detail-section="${escapeHTML(target)}">
      <span>${escapeHTML(label)}</span>
      <strong class="${className || ""}">${value}</strong>
      ${detail ? `<small>${escapeHTML(detail)}</small>` : ""}
    </button>
  `;
}

function stockDetailEvidenceCard(target, title, value, detail, tone = "") {
  return `
    <button class="stock-detail-evidence-card ${tone}" type="button" data-detail-section="${escapeHTML(target)}">
      <span>${escapeHTML(title)}</span>
      <strong>${value}</strong>
      <small>${escapeHTML(detail)}</small>
    </button>
  `;
}

function stockDetailAccordion(id, eyebrow, title, body) {
  const expanded = expandedStockDetailSections.has(id);
  return `
    <section class="stock-detail-section stock-detail-accordion ${expanded ? "" : "is-collapsed"}" id="${escapeHTML(id)}">
      <button class="stock-detail-accordion-head" type="button" data-stock-detail-toggle="${escapeHTML(id)}" aria-expanded="${expanded ? "true" : "false"}">
        <span>
          <small>${escapeHTML(eyebrow)}</small>
          <strong>${escapeHTML(title)}</strong>
        </span>
        <em>${expanded ? "收起" : "展开"}</em>
      </button>
      <div class="stock-detail-accordion-body">
        ${body}
      </div>
    </section>
  `;
}

function stockDetailInputValue(value) {
  return escapeHTML(value ?? "");
}

function stockDetailNumberValue(value, decimals = 2) {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return "";
  return escapeHTML(String(Number(number.toFixed(decimals))));
}

function stockDetailPercentInputValue(value, decimals = 2) {
  const number = finiteNumber(value);
  if (!Number.isFinite(number)) return "";
  return escapeHTML(String(Number((number * 100).toFixed(decimals))));
}

function stockDetailScenarioInputs(stock) {
  const scenarios = fallbackValuationScenarios(stock);
  return STOCK_DETAIL_VALUATION_SCENARIOS.map((meta, index) => {
    const scenario = scenarios[index] ?? {};
    return {
      ...scenario,
      key: meta.key,
      label: meta.label,
      name: meta.label
    };
  });
}

function stockDetailHumanInputPanel(stock) {
  return `
    <section class="stock-detail-section" id="detailInputs">
      <div class="stock-detail-section-head">
        <p class="eyebrow">Human Inputs</p>
        <h2>人工判断</h2>
      </div>
      <section class="panel stock-detail-edit-panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Editable</p>
            <h2>只编辑系统无法替你判断的内容</h2>
          </div>
        </div>
        <div class="detail-content">
          <form class="stock-detail-edit-form" data-stock-human-input-form data-symbol="${escapeHTML(stock.symbol)}">
            <div class="stock-detail-form-grid">
              <label class="wide">
                <span>买入逻辑</span>
                <textarea name="buyLogic" rows="3" placeholder="为什么值得买，哪些事实会让这个逻辑失效">${stockDetailInputValue(stock.buyLogic)}</textarea>
              </label>
              <label>
                <span>当前动作</span>
                <input name="action" type="text" value="${stockDetailInputValue(stock.action)}" placeholder="继续持有 / 等待 / 卖出" />
              </label>
              <label>
                <span>达标状态</span>
                <input name="status" type="text" value="${stockDetailInputValue(stock.status)}" placeholder="达标 / 观察 / 风险排除" />
              </label>
              <label>
                <span>估值可信度</span>
                <input name="valuationConfidence" type="text" value="${stockDetailInputValue(stock.valuationConfidence)}" placeholder="high / medium / low 或中文说明" />
              </label>
              <label>
                <span>质量总分</span>
                <input name="qualityScore" type="number" min="0" max="100" step="1" value="${stockDetailNumberValue(stock.qualityScore, 0)}" />
              </label>
              <label>
                <span>商业模式</span>
                <input name="businessModel" type="number" min="0" max="30" step="1" value="${stockDetailNumberValue(stock.businessModel, 0)}" />
              </label>
              <label>
                <span>护城河</span>
                <input name="moat" type="number" min="0" max="25" step="1" value="${stockDetailNumberValue(stock.moat, 0)}" />
              </label>
              <label>
                <span>治理</span>
                <input name="governance" type="number" min="0" max="20" step="1" value="${stockDetailNumberValue(stock.governance, 0)}" />
              </label>
              <label>
                <span>财务质量</span>
                <input name="financialQuality" type="number" min="0" max="25" step="1" value="${stockDetailNumberValue(stock.financialQuality, 0)}" />
              </label>
              <label class="wide">
                <span>主要风险</span>
                <textarea name="risk" rows="3" placeholder="最重要的业务、治理、财务或行业风险">${stockDetailInputValue(stock.risk)}</textarea>
              </label>
              <label class="wide">
                <span>反证条件</span>
                <textarea name="killCriteria" rows="3" placeholder="一行一个：出现什么事实就要推翻买入逻辑">${stockDetailInputValue(killCriteriaItems(stock).join("\n"))}</textarea>
              </label>
              <label class="wide">
                <span>备注</span>
                <textarea name="notes" rows="3" placeholder="补充判断，不放系统能自动获取的数据">${stockDetailInputValue(stock.notes)}</textarea>
              </label>
            </div>
            <div class="stock-detail-form-actions">
              <span>自动数据仍由行情、财报和估值分位更新维护。</span>
              <button class="primary-button compact-link" type="submit">保存人工判断</button>
            </div>
          </form>
        </div>
      </section>
      <section class="panel stock-detail-edit-panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Decision Log</p>
            <h2>继续持有 / 复盘理由</h2>
          </div>
        </div>
        <div class="detail-content">
          <form class="stock-detail-edit-form" data-stock-hold-log-form data-symbol="${escapeHTML(stock.symbol)}">
            <div class="stock-detail-form-grid">
              <label>
                <span>日志类型</span>
                <select name="decision">
                  <option value="继续持有">继续持有</option>
                  <option value="复盘">复盘</option>
                </select>
              </label>
              <label class="wide">
                <span>人工理由</span>
                <textarea name="detail" rows="3" required placeholder="写清楚这次为什么继续持有、复盘结论是什么"></textarea>
              </label>
            </div>
            <div class="stock-detail-form-actions">
              <span>保存时自动附带现价、仓位、估值区间和安全边际快照。</span>
              <button class="primary-button compact-link" type="submit">保存日志</button>
            </div>
          </form>
        </div>
      </section>
    </section>
  `;
}

function stockDetailValuationInputPanel(stock) {
  const range = valuationRangeView(stock);
  const requiredMargin = finiteNumber(stock?.valuation?.requiredMargin) ?? MAIN_DCF_MARGIN_TARGET;
  return `
    <section class="panel stock-detail-edit-panel">
      <div class="panel-head compact">
        <div>
          <p class="eyebrow">Editable Assumptions</p>
          <h2>三情景估值假设</h2>
        </div>
      </div>
      <div class="detail-content">
        <form class="stock-detail-edit-form" data-stock-valuation-form data-symbol="${escapeHTML(stock.symbol)}">
          <div class="stock-detail-valuation-topline">
            <label>
              <span>币种</span>
              <input name="currency" type="text" value="${stockDetailInputValue(range.currency || stock.currency || "CNY")}" />
            </label>
            <label>
              <span>当前价</span>
              <input name="currentPrice" type="number" step="0.0001" value="${stockDetailNumberValue(stock.currentPrice, 4)}" />
            </label>
            <label>
              <span>安全边际要求</span>
              <input name="requiredMargin" type="number" step="0.1" value="${stockDetailPercentInputValue(requiredMargin, 1)}" />
            </label>
            <div>
              <span>当前估值区间</span>
              <strong>${Number.isFinite(range.low) && Number.isFinite(range.high) ? `${currency(range.low, range.currency)} - ${currency(range.high, range.currency)}` : "待补充"}</strong>
              <small>${Number.isFinite(range.margin) ? `安全边际 ${percent(range.margin * 100, false)}` : range.source}</small>
            </div>
          </div>
          <div class="stock-detail-scenario-table" role="table" aria-label="三情景估值假设">
            <div class="stock-detail-scenario-row head" role="row">
              <span>情景</span>
              <span>收入增长%</span>
              <span>利润率%</span>
              <span>FCF</span>
              <span>折现率%</span>
              <span>PE</span>
              <span>P/FCF</span>
              <span>股本</span>
            </div>
            ${stockDetailScenarioInputs(stock).map((scenario) => `
              <div class="stock-detail-scenario-row" role="row">
                <strong>${escapeHTML(scenario.label)}</strong>
                <input name="${escapeHTML(scenario.key)}.revenueGrowth" type="number" step="0.1" value="${stockDetailPercentInputValue(scenario.revenueGrowth, 2)}" />
                <input name="${escapeHTML(scenario.key)}.profitMargin" type="number" step="0.1" value="${stockDetailPercentInputValue(scenario.profitMargin, 2)}" />
                <input name="${escapeHTML(scenario.key)}.fcf" type="number" step="0.01" value="${stockDetailNumberValue(scenario.fcf, 2)}" />
                <input name="${escapeHTML(scenario.key)}.discountRate" type="number" step="0.1" value="${stockDetailPercentInputValue(scenario.discountRate, 2)}" />
                <input name="${escapeHTML(scenario.key)}.reasonablePe" type="number" step="0.1" value="${stockDetailNumberValue(scenario.reasonablePe, 2)}" />
                <input name="${escapeHTML(scenario.key)}.reasonablePfcf" type="number" step="0.1" value="${stockDetailNumberValue(scenario.reasonablePFCF ?? scenario.reasonablePfcf, 2)}" />
                <input name="${escapeHTML(scenario.key)}.shares" type="number" step="0.01" value="${stockDetailNumberValue(scenario.shares, 2)}" />
              </div>
            `).join("")}
          </div>
          <div class="stock-detail-form-actions">
            <span>保存后写入估值区间，选股页和总览继续只读展示。</span>
            <button class="primary-button compact-link" type="submit">保存估值假设</button>
          </div>
        </form>
      </div>
    </section>
  `;
}

function stockDetailSummaryPanel(stock, strategy, health, plan) {
  return `
    <section class="stock-detail-section" id="detailSummary">
      <div class="stock-detail-section-head">
        <p class="eyebrow">Research Summary</p>
        <h2>研究摘要</h2>
      </div>
      <div class="stock-detail-summary-grid">
        <article>
          <span>当前动作</span>
          <strong>${escapeHTML(displayText(stock.action, "暂无动作"))}</strong>
        </article>
        <article>
          <span>达标状态</span>
          <strong>${escapeHTML(displayText(stock.status, "未填写"))}</strong>
        </article>
        <article>
          <span>策略归属</span>
          <strong>${escapeHTML(strategy.status)}</strong>
          <small>${escapeHTML(strategy.blockers.length ? strategy.blockers.join("；") : "当前无主要阻碍")}</small>
        </article>
        <article>
          <span>执行计划</span>
          <strong>${plan ? `${escapeHTML(String(plan.rank))}. ${escapeHTML(plan.priority)}` : "未列入优先级"}</strong>
          <small>${escapeHTML(plan ? plan.advice : displayText(stock.action, "暂无执行计划"))}</small>
        </article>
        <article>
          <span>健康状态</span>
          <strong>${health ? `<span class="health-pill ${health.tone}">${escapeHTML(health.status)}</span>` : "-"}</strong>
          <small>${escapeHTML(health ? `${health.score}分 · ${health.detail}` : "未持仓，仅保留研究信息")}</small>
        </article>
        <article>
          <span>主要风险</span>
          <strong>${escapeHTML(displayText(stock.risk, "未填写"))}</strong>
        </article>
      </div>
    </section>
  `;
}

function stockDetailEvidenceGrid(stock, strategy, health, confidence, detailHoldingValue) {
  const annual = financialAnnuals(stock);
  const latest = annual[0] ?? {};
  const audit = strategy.ownerAudit;
  const netCash = strategy.netCash;
  const dividendYield = calculatedDividendYield(stock);
  const shareholderReturnYield = calculatedShareholderReturnYield(stock);
  const qualityText = Number.isFinite(stock.qualityScore) ? `${stock.qualityScore}` : "-";
  const fcfRecord = positiveRecordRatio(stock, "freeCashFlow");
  const riskItems = killCriteriaItems(stock);

  return `
    <div class="stock-detail-evidence-grid">
      ${stockDetailEvidenceCard(
        "detailValuation",
        "估值",
        Number.isFinite(strategy.margin) ? percent(strategy.margin * 100, false) : "-",
        `主策略门槛 ${percent(MAIN_DCF_MARGIN_TARGET * 100, false)}；可信度 ${confidence.text}`,
        Number.isFinite(strategy.margin) && strategy.margin >= MAIN_DCF_MARGIN_TARGET ? "strong" : "watch"
      )}
      ${stockDetailEvidenceCard(
        "detailFinancials",
        "财务质量",
        `${qualityText} / ${audit.hasAudit ? audit.score : "-"}`,
        `ROE/ROIC ${financialRatio(latest.roe)} / ${financialRatio(latest.roic)}；FCF为正 ${Number.isFinite(fcfRecord) ? `${Math.round(fcfRecord * 100)}%` : "未知"}`,
        audit.tone
      )}
      ${stockDetailEvidenceCard(
        "detailIncome",
        "现金回报",
        privateText(`${displayDividendRatio(shareholderReturnYield)} / ${displayDividendRatio(strategy.shield.target)}`),
        `${dividendReliability(stock).text}；股息率 ${privateText(displayDividendRatio(dividendYield))}`,
        strategy.shield.passed ? "strong" : "wait"
      )}
      ${stockDetailEvidenceCard(
        "detailRisk",
        "风险反证",
        health ? `<span class="health-pill ${health.tone}">${escapeHTML(health.status)}</span>` : "-",
        riskItems.length ? riskItems[0] : "暂无明确反证条件",
        health?.tone || "watch"
      )}
      ${stockDetailEvidenceCard(
        "detailValuation",
        "持仓暴露",
        escapeHTML(detailHoldingValue),
        `现价 ${Number.isFinite(stock.currentPrice) ? currency(stock.currentPrice, stock.currency) : "-"}；公允区间 ${displayText(stock.fairValueRange)}`,
        "neutral"
      )}
      ${stockDetailEvidenceCard(
        "detailIncome",
        "烟蒂口径",
        financialMultiple(netCash.exCashPe),
        `ex-cash PE 门槛 ≤${strategy.peLimit}x；FCF yield ${financialRatio(netCash.fcfYield)}`,
        Number.isFinite(netCash.exCashPe) && netCash.exCashPe <= strategy.peLimit ? "strong" : "wait"
      )}
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
  const marginText = displayMarginOfSafety(stock);
  const qualityText = Number.isFinite(stock.qualityScore) ? `${stock.qualityScore}` : "-";
  const hasCurrentQuote = Number.isFinite(stock.currentPrice) && stock.currentPrice > 0;
  const hasPreviousQuote = Number.isFinite(stock.previousClose) && stock.previousClose > 0;
  const priceChange = hasCurrentQuote && hasPreviousQuote ? stock.currentPrice - stock.previousClose : null;
  const dayMetricClass = isHolding && Number.isFinite(stock.dayChange)
    ? stock.dayChange >= 0 ? "positive" : "negative"
    : Number.isFinite(priceChange)
      ? priceChange >= 0 ? "positive" : "negative"
      : "";
  const pnlClass = isHolding && Number.isFinite(stock.pnlCny)
    ? stock.pnlCny >= 0 ? "positive" : "negative"
    : "";
  const totalValue = positions.reduce((sum, item) => sum + item.marketValueCny, 0);
  const health = isHolding ? holdingHealth(stock, totalValue) : null;
  const confidence = confidenceMeta(stock);
  const strategy = strategyProfile(stock);
  const detailPnlMeta = isHolding && Math.abs(finiteNumber(stock.realizedPnlCny) ?? 0) >= 0.005 && Number.isFinite(stock.unrealizedPnlCny)
    ? privateText(`浮动 ${currency(stock.unrealizedPnlCny)} · 已实现 ${currency(stock.realizedPnlCny)}`)
    : isHolding ? privateText(percent(stock.pnlRate)) : "";
  const detailDayRate = isHolding && stock.marketValueCny
    ? privateText(percent((stock.dayChange / stock.marketValueCny) * 100))
    : Number.isFinite(priceChange) && stock.previousClose
      ? percent((priceChange / stock.previousClose) * 100)
      : "";
  const detailHoldingValue = isHolding ? privateText(currency(stock.marketValueCny)) : "-";
  const detailDayLabel = isHolding ? "今日盈亏" : "今日涨跌";
  const detailDayValue = isHolding
    ? privateText(currency(stock.dayChange))
    : Number.isFinite(priceChange) ? currency(priceChange, stock.currency) : "-";

  elements.stockDetail.innerHTML = `
    <div class="stock-detail-workbench">
      <aside class="stock-detail-summary">
        <section class="stock-detail-decision-hero">
          <div class="stock-detail-hero-actions">
            <a class="ghost-button detail-back" href="#holdings">返回持仓</a>
            <button class="ghost-button compact-link" type="button" data-update-financials="${escapeHTML(stock.symbol)}">
              <span>↻</span>
              更新财务
            </button>
          </div>
          <div>
            <p class="eyebrow">${escapeHTML(stock.symbol)} · ${escapeHTML(displayText(stock.industry, "未分类"))}</p>
            <h2>${escapeHTML(stock.name)}</h2>
            ${stockDetailCategoryControl(stock)}
          </div>
          <div class="stock-detail-hero-meta">
            <span>${escapeHTML(closeDateText(stock) || displayText(stock.updatedAt, "行情日期未知"))}</span>
            <small>${escapeHTML(stock.financials?.updatedAt ? `财务 ${stock.financials.updatedAt}` : "财务数据待更新")}</small>
          </div>
        </section>

        <section class="stock-detail-conclusion-card">
          <span>当前动作</span>
          <strong>${escapeHTML(displayText(stock.action, "暂无动作"))}</strong>
          <small>${escapeHTML(displayText(stock.status, strategy.status))}</small>
        </section>

        <section class="stock-detail-signal-grid">
          ${stockDetailMetricButton("detailValuation", "安全边际", escapeHTML(marginText), `目标 ${percent(MAIN_DCF_MARGIN_TARGET * 100, false)}`)}
          ${stockDetailMetricButton("detailIncome", "综合回报率", escapeHTML(privateText(`${displayDividendRatio(strategy.shield.value)} / ${displayDividendRatio(strategy.shield.target)}`)), privateText(strategy.shield.source))}
          ${stockDetailMetricButton("detailFinancials", "长期评分", escapeHTML(strategy.ownerAudit.hasAudit ? `${strategy.ownerAudit.score}/100` : "-"), strategy.ownerAudit.text)}
          ${stockDetailMetricButton("detailRisk", "健康评分", escapeHTML(health ? `${health.score}分` : "-"), health ? health.status : "未持仓")}
          ${stockDetailMetricButton("detailSummary", detailDayLabel, escapeHTML(detailDayValue), detailDayRate, privateClass(dayMetricClass))}
          ${stockDetailMetricButton("detailSummary", "累计盈亏", escapeHTML(isHolding ? privateText(currency(stock.pnlCny)) : "-"), detailPnlMeta, privateClass(pnlClass))}
        </section>
      </aside>

      <main class="stock-detail-evidence-flow">
        ${stockDetailNav()}

        ${stockDetailHumanInputPanel(stock)}

        ${stockDetailSummaryPanel(stock, strategy, health, plan)}

        <section class="stock-detail-section" id="detailValuation">
          <div class="stock-detail-section-head">
            <p class="eyebrow">Evidence Matrix</p>
            <h2>四象限证据</h2>
          </div>
          ${stockDetailValuationInputPanel(stock)}
          ${stockDetailEvidenceGrid(stock, strategy, health, confidence, detailHoldingValue)}

          <section class="detail-grid">
            <section class="panel">
              <div class="panel-head compact">
                <div>
                  <p class="eyebrow">Analysis</p>
                  <h2>研究判断</h2>
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
                  <div><span>综合回报率</span><strong>${escapeHTML(privateText(`${displayDividendRatio(strategy.shield.value)} / ${displayDividendRatio(strategy.shield.target)}`))}</strong><small>${escapeHTML(privateText(strategy.shield.source))}</small></div>
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
                  <div><span>持仓市值</span><strong>${escapeHTML(detailHoldingValue)}</strong></div>
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

          ${renderMasterVotesPanel(stock, totalValue)}
          ${renderOwnerAuditPanel(stock)}
        </section>

        ${stockDetailAccordion("detailFinancials", "Financials", "财务质量", renderFinancialsPanel(stock, { collapsibleTable: true }))}

        ${stockDetailAccordion("detailIncome", "Income", "现金回报", `
          ${renderDividendPanel(stock, isHolding)}
          ${renderNetCashPanel(stock)}
          ${renderDataSourcePanel(stock)}
        `)}

        ${stockDetailAccordion("detailRisk", "Risk", "风险反证", `
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
        `)}

        ${stockDetailAccordion("detailRecords", "Records", "日志档案", `
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
        `)}
      </main>
    </div>
  `;
  requestStockDetailActiveUpdate();
}

function renderTradeStockNames() {
  if (!elements.tradeStockNames) return;
  const seen = new Set();
  const items = [
    ...(state.holdings ?? []).map((stock) => ({ ...stock, optionType: "股票" })),
    ...(state.candidates ?? []).map((stock) => ({ ...stock, optionType: "跟踪" }))
  ]
    .filter((item) => item?.name && item?.symbol)
    .filter((item) => {
      const key = `stock:${normalizeSymbol(item.symbol)}`;
      if (!key || seen.has(key)) return false;
      seen.add(key);
      return true;
    });
  elements.tradeStockNames.innerHTML = items
    .map((item) => `<option value="${escapeHTML(item.name)}" label="${escapeHTML(`${item.optionType} · ${item.symbol}`)}"></option>`)
    .join("");
}

function routeInfo(rawHash = window.location.hash.slice(1)) {
  const view = rawHash || "overview";
  if (view === "screener" || view === "valuation" || view === "sunny30" || view === "candidates") {
    return { view: "screener", page: "screener" };
  }
  if (view === "holdings" || view === "masters" || view === "positions") {
    return { view: "holdings", page: "holdings" };
  }
  if (view === "logs") {
    return { view, page: "trades" };
  }
  if (view.startsWith("industry=")) {
    return { view, page: "industry-detail", id: decodeURIComponent(view.slice("industry=".length)) };
  }
  if (view.startsWith("stock=")) {
    return { view, page: "stock-detail", id: decodeURIComponent(view.slice("stock=".length)) };
  }
  return { view, page: pageTitles[view] ? view : "overview" };
}

function renderContext() {
  const positions = computePositions();
  return { positions };
}

function renderRecordCount() {
  if (!elements.recordCount) return;
  elements.recordCount.textContent = privateText(`${state.stocks?.length ?? state.holdings.length + state.candidates.length} 只股票 · ${state.trades.length} 条交易`);
}

function renderLoadingState() {
  setQuoteUpdateStatus("正在加载组合数据...");
  [
    elements.totalAssetsMetric,
    elements.totalValue,
    elements.totalPositionPnl,
    elements.dayChange,
    elements.annualDividend,
    elements.dataQualityMetric
  ].filter(Boolean).forEach((element) => {
    element.textContent = "加载中";
    element.className = "";
  });
  if (elements.positionCount) elements.positionCount.textContent = "等待后端数据";
  if (elements.totalPositionPnlRate) elements.totalPositionPnlRate.textContent = "";
  if (elements.dayChangeRate) elements.dayChangeRate.textContent = "";
  if (elements.portfolioDividendYield) elements.portfolioDividendYield.textContent = "";
  if (elements.dataQualityDetail) elements.dataQualityDetail.textContent = "正在读取数据目录";
  if (elements.committeeConsensus) {
    elements.committeeConsensus.innerHTML = `<div class="empty-state compact-empty">正在加载组合数据...</div>`;
  }
}

function render(rawHash = window.location.hash.slice(1)) {
  const route = routeInfo(rawHash);
  const { positions } = renderContext();

  renderTradeStockNames();
  renderQuoteUpdateStatus(positions);
  renderRecordCount();

  if (route.page === "overview") {
    renderMetrics(positions);
    renderAssetAllocation(positions);
    renderDecisionArea(positions);
    renderCommitteeOverview(positions);
    return;
  }
  if (route.page === "holdings") {
    renderPositions(positions);
    renderMastersPage(positions);
    renderPlanAndCandidates();
    return;
  }
  if (route.page === "screener") {
    renderSunny30(positions);
    renderValuationModule(positions);
    return;
  }
  if (route.page === "industry") {
    renderIndustryDesk(positions);
    return;
  }
  if (route.page === "industry-detail") {
    renderIndustryDetail(positions, route.id);
    return;
  }
  if (route.page === "trades") {
    renderTrades();
    renderDecisionLogs();
    return;
  }
  if (route.page === "stock-detail") {
    renderStockDetail(positions, route.id);
  }
}

function findTradeStock(input) {
  const text = String(input ?? "").trim();
  if (!text) return null;
  const normalized = normalizeSymbol(text);
  const stocks = [...(state.holdings ?? []), ...(state.candidates ?? [])];
  return (
    stocks.find((stock) => normalizeSymbol(stock.symbol) === normalized) ||
    stocks.find((stock) => String(stock.name ?? "").trim() === text) ||
    stocks.find((stock) => {
      const name = String(stock.name ?? "").trim();
      return name && (name.includes(text) || text.includes(name));
    }) ||
    null
  );
}

function inferTradeCurrency(stock) {
  const currencyCode = String(stock?.currency ?? "").trim().toUpperCase();
  if (currencyCode) return currencyCode;
  const symbol = normalizeSymbol(stock?.symbol);
  if (symbol.endsWith(".HK")) return "HKD";
  if (symbol.endsWith(".SH") || symbol.endsWith(".SZ")) return "CNY";
  return "CNY";
}

function candidateFromHolding(holding) {
  const { shares, cost, ...candidate } = holding;
  return candidate;
}

function clearedCandidateFromHolding(holding) {
  const candidate = candidateFromHolding(holding);
  const status = String(candidate.status ?? "").trim();
  const action = String(candidate.action ?? "").trim();
  if (!status || status.includes("持仓")) {
    candidate.status = "晴仓30跟踪（清仓后）";
  }
  if (!action) {
    candidate.action = "清仓后继续放在晴仓30跟踪；等待重新达到买入纪律";
  } else if (action.includes("继续持有")) {
    candidate.action = action.replace("继续持有", "清仓后继续晴仓30跟踪");
  } else if (!action.includes("清仓后")) {
    candidate.action = `清仓后继续晴仓30跟踪；${action}`;
  }
  return candidate;
}

function upsertCandidate(candidate) {
  const symbol = normalizeSymbol(candidate?.symbol);
  const index = state.candidates.findIndex((item) => normalizeSymbol(item.symbol) === symbol);
  if (index >= 0) {
    state.candidates[index] = candidate;
  } else {
    state.candidates.push(candidate);
  }
}

function removeCandidate(symbol) {
  const normalized = normalizeSymbol(symbol);
  state.candidates = state.candidates.filter((item) => normalizeSymbol(item.symbol) !== normalized);
}

function optionalFormNumber(formData, name) {
  return finiteNumber(formData.get(name));
}

function stockForDetailSave(symbol) {
  const normalized = normalizeSymbol(symbol);
  return (
    (state.stocks ?? []).find((item) => normalizeSymbol(item.symbol) === normalized) ||
    (state.holdings ?? []).find((item) => normalizeSymbol(item.symbol) === normalized) ||
    (state.candidates ?? []).find((item) => normalizeSymbol(item.symbol) === normalized) ||
    null
  );
}

function optionalDetailFormNumber(formData, name) {
  const raw = String(formData.get(name) ?? "").trim();
  if (!raw) return null;
  const value = Number(raw);
  if (!Number.isFinite(value)) throw new Error(`${name} 必须是数字`);
  return value;
}

function optionalDetailFormRate(formData, name) {
  const value = optionalDetailFormNumber(formData, name);
  return value === null ? null : value / 100;
}

function inferStockCurrency(symbol) {
  const normalized = normalizeSymbol(symbol);
  if (normalized.endsWith(".HK")) return "HKD";
  if (normalized.endsWith(".US")) return "USD";
  return "CNY";
}

function localUpsertStockPayload(payload) {
  const normalized = normalizeSymbol(payload.symbol);
  state.stocks ??= [];
  const next = stockPayloadFromLegacy(payload);
  const stockIndex = state.stocks.findIndex((item) => normalizeSymbol(item.symbol) === normalized);
  if (stockIndex >= 0) state.stocks[stockIndex] = { ...state.stocks[stockIndex], ...next };
  else state.stocks.push(next);
  state.holdings = stocksToHoldings(state.stocks);
  state.candidates = stocksToCandidates(state.stocks);
  saveState();
  render();
}

async function saveStockPayload(symbol, patch) {
  const normalized = normalizeSymbol(symbol);
  const existing = stockForDetailSave(normalized);
  if (!existing) throw new Error(`未找到 ${normalized}`);
  const payload = stockPayloadFromLegacy({
    ...existing,
    ...patch,
    symbol: normalized,
    name: existing.name || patch.name || normalized,
    currency: patch.currency || existing.currency || inferStockCurrency(normalized)
  });

  if (USE_BACKEND) {
    setLoadedState(await requestJSON(`/api/stocks/${encodeURIComponent(normalized)}`, {
      method: "PUT",
      body: JSON.stringify(payload)
    }));
    localStorage.removeItem(STORAGE_KEY);
    render();
    return;
  }

  localUpsertStockPayload(payload);
}

async function saveStockHumanInputs(form) {
  const formData = new FormData(form);
  const symbol = normalizeSymbol(form.dataset.symbol);
  const killCriteria = String(formData.get("killCriteria") ?? "")
    .split(/\n+/)
    .map((item) => item.trim())
    .filter(Boolean);
  await saveStockPayload(symbol, {
    buyLogic: String(formData.get("buyLogic") ?? "").trim(),
    action: String(formData.get("action") ?? "").trim(),
    status: String(formData.get("status") ?? "").trim(),
    valuationConfidence: String(formData.get("valuationConfidence") ?? "").trim(),
    risk: String(formData.get("risk") ?? "").trim(),
    notes: String(formData.get("notes") ?? "").trim(),
    killCriteria,
    qualityScore: optionalDetailFormNumber(formData, "qualityScore"),
    businessModel: optionalDetailFormNumber(formData, "businessModel"),
    moat: optionalDetailFormNumber(formData, "moat"),
    governance: optionalDetailFormNumber(formData, "governance"),
    financialQuality: optionalDetailFormNumber(formData, "financialQuality")
  });
}

function scenarioFairValueFromInputs(scenario) {
  const values = [];
  if (Number.isFinite(scenario.fairValue) && scenario.fairValue > 0) values.push(scenario.fairValue);
  if (Number.isFinite(scenario.fcf) && scenario.fcf > 0 && Number.isFinite(scenario.shares) && scenario.shares > 0) {
    if (Number.isFinite(scenario.reasonablePfcf) && scenario.reasonablePfcf > 0) {
      values.push((scenario.fcf * scenario.reasonablePfcf) / scenario.shares);
    }
    if (Number.isFinite(scenario.reasonablePe) && scenario.reasonablePe > 0) {
      values.push((scenario.fcf * scenario.reasonablePe) / scenario.shares);
    }
  }
  if (!values.length) return null;
  return values.reduce((sum, value) => sum + value, 0) / values.length;
}

function valuationRangeFromScenarios(scenarios, currentPrice, currencyCode) {
  const values = scenarios
    .map((scenario) => scenarioFairValueFromInputs(scenario))
    .filter((value) => Number.isFinite(value) && value > 0)
    .sort((a, b) => a - b);
  if (!values.length) throw new Error("至少填写一个可计算的估值情景：FCF、倍数、股本");
  const base = values[Math.floor(values.length / 2)];
  const margin = Number.isFinite(currentPrice) && currentPrice > 0 && base > 0 ? (base - currentPrice) / base : null;
  return {
    low: values[0],
    base,
    high: values[values.length - 1],
    currency: currencyCode,
    marginOfSafety: margin
  };
}

async function saveStockValuationInputs(form) {
  const formData = new FormData(form);
  const symbol = normalizeSymbol(form.dataset.symbol);
  const existing = stockForDetailSave(symbol);
  if (!existing) throw new Error(`未找到 ${symbol}`);
  const currencyCode = String(formData.get("currency") || existing.currency || inferStockCurrency(symbol)).trim().toUpperCase();
  const currentPrice = optionalDetailFormNumber(formData, "currentPrice") ?? finiteNumber(existing.currentPrice) ?? 0;
  const requiredMargin = optionalDetailFormRate(formData, "requiredMargin") ?? MAIN_DCF_MARGIN_TARGET;
  const scenarios = STOCK_DETAIL_VALUATION_SCENARIOS.map((meta) => ({
    name: meta.label,
    revenueGrowth: optionalDetailFormRate(formData, `${meta.key}.revenueGrowth`),
    profitMargin: optionalDetailFormRate(formData, `${meta.key}.profitMargin`),
    fcf: optionalDetailFormNumber(formData, `${meta.key}.fcf`),
    discountRate: optionalDetailFormRate(formData, `${meta.key}.discountRate`),
    reasonablePe: optionalDetailFormNumber(formData, `${meta.key}.reasonablePe`),
    reasonablePfcf: optionalDetailFormNumber(formData, `${meta.key}.reasonablePfcf`),
    shares: optionalDetailFormNumber(formData, `${meta.key}.shares`)
  }));
  const range = valuationRangeFromScenarios(scenarios, currentPrice, currencyCode);
  await saveStockPayload(symbol, {
    currentPrice,
    currency: currencyCode,
    valuation: {
      currency: currencyCode,
      currentPrice,
      requiredMargin,
      updatedAt: new Date().toISOString().slice(0, 10),
      source: "人工详情页假设",
      scenarios,
      range
    },
    intrinsicValue: range.base,
    fairValueRange: `${currency(range.low, currencyCode)} - ${currency(range.high, currencyCode)}`,
    marginOfSafety: range.marginOfSafety
  });
}

async function saveStockHoldLog(form) {
  const formData = new FormData(form);
  const symbol = normalizeSymbol(form.dataset.symbol);
  const detail = String(formData.get("detail") ?? "").trim();
  if (!detail) throw new Error("日志理由不能为空");
  const decision = String(formData.get("decision") || "继续持有").trim();
  if (USE_BACKEND) {
    setLoadedState(await requestJSON("/api/decision-logs", {
      method: "POST",
      body: JSON.stringify({
        type: decision === "复盘" ? "review" : "hold",
        symbol,
        decision,
        detail
      })
    }));
    localStorage.removeItem(STORAGE_KEY);
    render();
    return;
  }
  state.decisionLogs ??= [];
  state.decisionLogs.push({
    id: Date.now(),
    date: new Date().toISOString(),
    type: decision === "复盘" ? "review" : "hold",
    symbol,
    decision,
    detail
  });
  saveState();
  render();
}

function candidateFromSunny30Form(formData) {
  const symbol = normalizeSymbol(formData.get("symbol"));
  const existingHolding = (state.holdings ?? []).find((holding) => normalizeSymbol(holding.symbol) === symbol);
  const name = String(formData.get("name") ?? "").trim() || existingHolding?.name || "";
  const category = String(formData.get("category") ?? "修复仓").trim() || "修复仓";
  const currencyCode = String(formData.get("currency") ?? "").trim().toUpperCase() || inferTradeCurrency({ symbol });
  const currentPrice = optionalFormNumber(formData, "currentPrice");
  const intrinsicValue = optionalFormNumber(formData, "intrinsicValue");
  const qualityScore = optionalFormNumber(formData, "qualityScore");
  const moat = optionalFormNumber(formData, "moat");

  if (!symbol) throw new Error("请填写股票代码");
  if (!name) throw new Error("请填写股票名称");
  if (currentPrice !== null && currentPrice <= 0) throw new Error("最新价必须大于 0");
  if (intrinsicValue !== null && intrinsicValue <= 0) throw new Error("内在价值必须大于 0");
  if (qualityScore !== null && (qualityScore < 0 || qualityScore > 100)) throw new Error("公司质量需在 0-100 之间");
  if (moat !== null && (moat < 0 || moat > 25)) throw new Error("护城河需在 0-25 之间");

  const today = new Date().toISOString().slice(0, 10);
  const candidate = {
    symbol,
    name,
    category,
    industry: existingHolding?.industry || category,
    currency: currencyCode,
    updatedAt: today,
    notes: String(formData.get("notes") ?? "").trim()
  };
  if (currentPrice !== null) {
    candidate.currentPrice = currentPrice;
    candidate.previousClose = currentPrice;
    candidate.currentPriceDate = today;
    candidate.previousCloseDate = today;
  }
  if (intrinsicValue !== null) candidate.intrinsicValue = intrinsicValue;
  if (qualityScore !== null) candidate.qualityScore = qualityScore;
  if (moat !== null) candidate.moat = moat;
  const margin = calculatedMarginOfSafety(candidate);
  if (Number.isFinite(margin)) candidate.marginOfSafety = margin;

  return candidate;
}

function openSunny30CandidateDialog() {
  if (!elements.sunny30CandidateDialog || !elements.sunny30CandidateForm) return;
  elements.sunny30CandidateForm.reset();
  elements.sunny30CandidateForm.category.value = "修复仓";
  elements.sunny30CandidateDialog.showModal();
}

async function saveSunny30Candidate(formData) {
  const candidate = candidateFromSunny30Form(formData);

  if (USE_BACKEND) {
    setLoadedState(await requestJSON("/api/stocks", {
      method: "POST",
      body: JSON.stringify(candidate)
    }));
    localStorage.removeItem(STORAGE_KEY);
    render("screener");
    return;
  }

  state.candidates ??= [];
  const existing = state.candidates.find((item) => normalizeSymbol(item.symbol) === candidate.symbol);
  const nextCandidate = { ...(existing ?? {}), ...candidate };
  nextCandidate.status ||= "晴仓30跟踪";
  nextCandidate.action ||= "纳入晴仓30长期跟踪；等待质量和安全边际补充";
  upsertCandidate(nextCandidate);
  saveState();
  render("screener");
}

async function deleteSunny30Candidate(symbol) {
  const normalized = normalizeSymbol(symbol);
  if (!normalized) return;
  const existing = sunny30Universe(computePositions()).find((stock) => normalizeSymbol(stock.symbol) === normalized);
  if (existing && !sunny30CanDelete(existing)) {
    throw new Error("持仓标的不能从选股页删除，请先在持仓页处理");
  }

  if (USE_BACKEND) {
    setLoadedState(await requestJSON(`/api/stocks/${encodeURIComponent(normalized)}`, { method: "DELETE" }));
    localStorage.removeItem(STORAGE_KEY);
    render("screener");
    return;
  }

  removeCandidate(normalized);
  saveState();
  render("screener");
}

async function saveScreeningWeights(formData) {
  const weights = SCREENING_WEIGHT_FIELDS.reduce((next, field) => {
    next[field.key] = Number(formData.get(field.key));
    if (!Number.isFinite(next[field.key]) || next[field.key] < 0) throw new Error(`${field.label}权重无效`);
    return next;
  }, {});
  const total = SCREENING_WEIGHT_FIELDS.reduce((sum, field) => sum + weights[field.key], 0);
  if (total !== 100) throw new Error("筛选权重合计必须为 100");

  if (USE_BACKEND) {
    setLoadedState(await requestJSON("/api/screening-weights", {
      method: "PUT",
      body: JSON.stringify(weights)
    }));
    localStorage.removeItem(STORAGE_KEY);
    render("screener");
    return;
  }

  state.screeningWeights = weights;
  saveState();
  render("screener");
}

function categoryLabelFromValue(value) {
  const key = normalizePositionCategory(value);
  const category = POSITION_CATEGORY_META[key];
  if (!category) throw new Error("请选择有效分类");
  return category.label;
}

async function saveStockCategory(symbol, value) {
  const normalized = normalizeSymbol(symbol);
  const category = categoryLabelFromValue(value);
  const holding = (state.holdings ?? []).find((item) => normalizeSymbol(item.symbol) === normalized);
  const candidate = (state.candidates ?? []).find((item) => normalizeSymbol(item.symbol) === normalized);
  if (!holding && !candidate) throw new Error(`未找到 ${normalized}`);

  if (USE_BACKEND) {
    const payload = stockPayloadFromLegacy({
      ...(candidate ?? {}),
      ...(holding ?? {}),
      symbol: normalized,
      name: candidate?.name || holding?.name || normalized,
      category
    });
    setLoadedState(await requestJSON(`/api/stocks/${encodeURIComponent(normalized)}`, {
      method: "PUT",
      body: JSON.stringify(payload)
    }));
    localStorage.removeItem(STORAGE_KEY);
    render();
    return category;
  }

  if (holding) holding.category = category;
  if (candidate) candidate.category = category;
  saveState();
  render();
  return category;
}

function tradeFromSimpleForm(formData) {
  const nameInput = String(formData.get("name") ?? "").trim();
  const signedShares = Number(formData.get("shares"));
  const price = Number(formData.get("price"));
  const reason = String(formData.get("reason") ?? "").trim();
  if (!nameInput) throw new Error("请填写股票名称");
  if (!Number.isFinite(signedShares) || signedShares === 0) throw new Error("股数不能为 0；买入填正数，卖出填负数");
  if (!Number.isFinite(price) || price <= 0) throw new Error("成交价必须大于 0");
  if (!reason) throw new Error("请填写买入、卖出或继续持有的人工理由");

  const stock = findTradeStock(nameInput);
  if (!stock) throw new Error(`未找到“${nameInput}”，请先把它加入持仓或晴仓30`);

  const symbol = normalizeSymbol(stock.symbol);
  const currencyCode = inferTradeCurrency(stock);
  const currentPrice = Number(stock.currentPrice) > 0 ? Number(stock.currentPrice) : price;
  const side = signedShares < 0 ? "sell" : "buy";
  const shares = Math.abs(signedShares);
  return {
    id: Date.now(),
    date: new Date().toISOString().slice(0, 10),
    assetType: "stock",
    symbol,
    name: stock.name || nameInput,
    side,
    shares,
    price,
    currency: currencyCode,
    currentPrice,
    reason
  };
}

async function addTrade(formData) {
  const trade = tradeFromSimpleForm(formData);
  const { side, shares, price, symbol } = trade;
  const currencyCode = trade.currency;

  if (USE_BACKEND) {
    setLoadedState(await requestJSON("/api/trades", {
      method: "POST",
      body: JSON.stringify(trade)
    }));
    saveState();
    render();
    return;
  }

  let holding = state.holdings.find((item) => normalizeSymbol(item.symbol) === symbol);
  if (!holding) {
    if (side === "sell") throw new Error("未找到可卖出的持仓");
    const candidate = state.candidates.find((item) => normalizeSymbol(item.symbol) === symbol);
    holding = candidate
      ? { ...candidate, shares: 0, cost: price }
      : {
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
    removeCandidate(symbol);
    state.holdings.push(holding);
  }

  if (side === "buy") {
    removeCandidate(symbol);
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
    detail: `${sideText} ${shares} 股；成交价 ${currencyCode} ${price.toFixed(4)}；录入最新价 ${currencyCode} ${trade.currentPrice.toFixed(4)}；理由：${trade.reason}`
  });
  if (side === "sell" && holding.shares === 0) {
    upsertCandidate(clearedCandidateFromHolding(holding));
    state.holdings = state.holdings.filter((item) => normalizeSymbol(item.symbol) !== symbol);
  }
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

  setLoadedState(result.state);
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
  setQuoteUpdateStatus("正在拉取股票行情和股息数据...");

  try {
    const result = await requestJSON("/api/quotes/update", { method: "POST" });
    setLoadedState(result.state);
    localStorage.removeItem(STORAGE_KEY);
    syncCash();
    render();

    const skipped = result.skipped ?? [];
    if (skipped.length) {
      const preview = skipped.slice(0, 3).map((item) => `${item.symbol} ${item.error}`).join("；");
      setQuoteUpdateStatus(`已更新 ${result.updated} 项行情，${skipped.length} 个失败：${preview}`, "error");
    } else {
      setQuoteUpdateStatus(`已更新 ${result.updated} 项行情`, "success");
    }
  } finally {
    elements.updateQuotesButton.disabled = false;
    elements.updateQuotesButton.innerHTML = "<span>↻</span> 更新行情";
  }
}

async function updateValuationHistory() {
  if (!USE_BACKEND) throw new Error("需要通过 go run . 启动后端后才能更新估值分位");
  const button = elements.updateValuationHistoryButton;
  const originalHTML = button?.innerHTML;
  if (button) {
    button.disabled = true;
    button.innerHTML = "<span>↻</span> 更新中";
  }
  setQuoteUpdateStatus("正在更新估值历史分位...");
  try {
    const result = await requestJSON("/api/valuation-history/update", { method: "POST" });
    if (result.state) {
      setLoadedState(result.state);
      localStorage.removeItem(STORAGE_KEY);
      render("valuation");
    }
    setQuoteUpdateStatus(`已写入 ${result.updated ?? 0} 只股票的估值分位历史`, "success");
  } finally {
    if (button) {
      button.disabled = false;
      button.innerHTML = originalHTML;
    }
  }
}

async function updateIndustryMetrics(button) {
  if (!USE_BACKEND) throw new Error("需要通过 go run . 启动后端后才能更新行业数据");

  const originalHTML = button?.innerHTML;
  if (button) {
    button.disabled = true;
    button.innerHTML = "<span>↻</span> 更新中";
  }
  setQuoteUpdateStatus("正在更新行业外部趋势数据...");

  try {
    const result = await requestJSON("/api/industries/update", { method: "POST" });
    setLoadedState(result.state);
    localStorage.removeItem(STORAGE_KEY);
    render();

    const skipped = result.skipped ?? [];
    if (skipped.length) {
      const preview = skipped.slice(0, 2).map((item) => `${item.industryId || item.source} ${item.error}`).join("；");
      setQuoteUpdateStatus(`已更新 ${result.updated} 个行业趋势指标，${skipped.length} 个来源失败：${preview}`, "error");
    } else {
      setQuoteUpdateStatus(`已更新 ${result.updated} 个行业趋势指标`, "success");
    }
  } finally {
    if (button) {
      button.disabled = false;
      button.innerHTML = originalHTML;
    }
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
    setLoadedState(result.state);
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

function openHoldingEditor(symbol) {
  const holding = state.holdings.find((item) => item.symbol === symbol);
  if (!holding) return;

  const form = elements.holdingForm;
  form.symbol.value = holding.symbol;
  form.name.value = holding.name ?? "";
  form.industry.value = holding.industry ?? "";
  form.category.value = positionCategory(holding).label;
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
    category: String(formData.get("category") ?? "").trim(),
    action: String(formData.get("action")).trim(),
    status: String(formData.get("status")).trim(),
    marginOfSafety: calculatedMarginOfSafety(holding),
    qualityScore: formData.get("qualityScore") === "" ? null : Number(formData.get("qualityScore")),
    notes: String(formData.get("notes")).trim()
  };

  if (USE_BACKEND) {
    setLoadedState(await requestJSON(`/api/stocks/${encodeURIComponent(symbol)}`, {
      method: "PUT",
      body: JSON.stringify(stockPayloadFromLegacy({ ...holding, ...patch }))
    }));
    saveState();
    render();
    return;
  }

  holding.name = patch.name;
  holding.industry = patch.industry;
  holding.category = patch.category;
  holding.action = patch.action;
  holding.status = patch.status;
  holding.marginOfSafety = patch.marginOfSafety;
  holding.qualityScore = patch.qualityScore;
  holding.notes = patch.notes;
  saveState();
  render();
}

function setTradeAssetType() {
  if (elements.tradeNameLabel) elements.tradeNameLabel.textContent = "股票名称";
  if (elements.tradePriceLabel) elements.tradePriceLabel.textContent = "成交价";
  if (elements.tradeSharesLabel) elements.tradeSharesLabel.textContent = "股数";
  if (elements.tradeNameInput) elements.tradeNameInput.placeholder = "中海物业";
  if (elements.tradeSharesInput) {
    elements.tradeSharesInput.placeholder = "买入填正数，卖出填负数";
    elements.tradeSharesInput.step = "1";
  }
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

let backToTopTicking = false;
let stockDetailActiveTicking = false;

function updateBackToTopVisibility() {
  elements.backToTopButton?.classList.toggle("is-visible", window.scrollY > 360);
}

function requestBackToTopUpdate() {
  if (backToTopTicking) return;
  backToTopTicking = true;
  requestAnimationFrame(() => {
    backToTopTicking = false;
    updateBackToTopVisibility();
  });
}

window.addEventListener("scroll", requestBackToTopUpdate, { passive: true });

elements.backToTopButton?.addEventListener("click", () => {
  window.scrollTo({ top: 0, behavior: "smooth" });
});

document.addEventListener("click", (event) => {
  const stockLink = event.target.closest('a[href^="#stock="]');
  if (!stockLink) return;
  savePortfolioReturnScroll();
}, { capture: true });

function setStockDetailActiveSection(sectionId) {
  if (!sectionId) return;
  activeStockDetailSection = sectionId;
  document.querySelectorAll(".stock-detail-nav [data-detail-section]").forEach((button) => {
    const active = button.dataset.detailSection === sectionId;
    button.classList.toggle("active", active);
    button.setAttribute("aria-pressed", active ? "true" : "false");
  });
}

function updateStockDetailToggle(section, expanded) {
  const toggle = section?.querySelector("[data-stock-detail-toggle]");
  if (!toggle) return;
  toggle.setAttribute("aria-expanded", expanded ? "true" : "false");
  const label = toggle.querySelector("em");
  if (label) label.textContent = expanded ? "收起" : "展开";
}

function expandStockDetailSection(sectionId) {
  const section = document.getElementById(sectionId);
  if (!section?.classList.contains("stock-detail-accordion")) return;
  expandedStockDetailSections.add(sectionId);
  section.classList.remove("is-collapsed");
  updateStockDetailToggle(section, true);
}

function toggleStockDetailSection(sectionId) {
  const section = document.getElementById(sectionId);
  if (!section?.classList.contains("stock-detail-accordion")) return;
  const expanded = !expandedStockDetailSections.has(sectionId);
  if (expanded) {
    expandedStockDetailSections.add(sectionId);
  } else {
    expandedStockDetailSections.delete(sectionId);
  }
  section.classList.toggle("is-collapsed", !expanded);
  updateStockDetailToggle(section, expanded);
}

function updateStockDetailActiveFromScroll() {
  const activePage = document.querySelector('.page[data-page="stock-detail"].active');
  if (!activePage) return;
  const sections = STOCK_DETAIL_NAV_ITEMS
    .map((item) => document.getElementById(item.id))
    .filter(Boolean);
  if (!sections.length) return;

  const anchor = window.innerWidth <= 720 ? 96 : 88;
  const viewportLimit = window.innerHeight * 0.72;
  let current = sections[0].id;
  let closestDistance = Number.POSITIVE_INFINITY;
  sections.forEach((section) => {
    const top = section.getBoundingClientRect().top;
    if (top > viewportLimit) return;
    const distance = Math.abs(top - anchor);
    if (distance < closestDistance) {
      closestDistance = distance;
      current = section.id;
    }
  });
  setStockDetailActiveSection(current);
}

function requestStockDetailActiveUpdate() {
  if (stockDetailActiveTicking) return;
  stockDetailActiveTicking = true;
  requestAnimationFrame(() => {
    stockDetailActiveTicking = false;
    updateStockDetailActiveFromScroll();
  });
}

window.addEventListener("scroll", requestStockDetailActiveUpdate, { passive: true });

function showPage(view) {
  const route = routeInfo(view);
  const isStockDetail = route.page === "stock-detail";
  const isIndustryDetail = route.page === "industry-detail";
  let nextView = route.page;
  if (!document.querySelector(`[data-page="${nextView}"]`)) {
    nextView = "overview";
  }

  document.querySelectorAll(".nav-item.active").forEach((item) => item.classList.remove("active"));
  const activeNavView = isIndustryDetail ? "industry" : isStockDetail ? "holdings" : route.view;
  document.querySelectorAll(`.nav-item[data-view="${activeNavView}"]`).forEach((item) => item.classList.add("active"));
  document.querySelector(".page.active")?.classList.remove("active");
  document.querySelector(`[data-page="${nextView}"]`)?.classList.add("active");
  elements.pageTitle.textContent = pageTitles[route.view] || pageTitles[nextView];
  requestBackToTopUpdate();
}

function resetStockDetailRouteScroll() {
  activeStockDetailSection = "detailInputs";
  requestAnimationFrame(() => {
    window.scrollTo({ top: 0, left: 0, behavior: "auto" });
    setStockDetailActiveSection("detailInputs");
    requestBackToTopUpdate();
  });
}

function savePortfolioReturnScroll() {
  const activePortfolio = document.querySelector('.page[data-page="holdings"].active, .page[data-page="screener"].active');
  if (!activePortfolio) return;
  sessionStorage.setItem(PORTFOLIO_RETURN_SCROLL_KEY, String(Math.max(0, Math.round(window.scrollY))));
}

function restorePortfolioReturnScroll() {
  const stored = Number(sessionStorage.getItem(PORTFOLIO_RETURN_SCROLL_KEY));
  sessionStorage.removeItem(PORTFOLIO_RETURN_SCROLL_KEY);
  if (!Number.isFinite(stored) || stored <= 0) return;
  requestAnimationFrame(() => {
    window.scrollTo({ top: stored, left: 0, behavior: "auto" });
    requestBackToTopUpdate();
  });
}

function handleRoute(rawHash) {
  const view = rawHash || "overview";
  const route = routeInfo(view);
  const previousRoute = activeRoute;
  showPage(view);
  render(view);
  if (route.page === "stock-detail") {
    resetStockDetailRouteScroll();
  } else if ((route.page === "holdings" || route.page === "screener") && previousRoute?.page === "stock-detail") {
    restorePortfolioReturnScroll();
  }
  activeRoute = route;
}

window.addEventListener("hashchange", () => {
  handleRoute(window.location.hash.slice(1));
});

elements.privacyToggle?.addEventListener("click", () => {
  holdingsMasked = !holdingsMasked;
  localStorage.setItem(HOLDINGS_PRIVACY_KEY, holdingsMasked ? "1" : "0");
  setPrivacyToggleState();
  render();
});

elements.candidateSort?.addEventListener("change", (event) => {
  candidateSort = event.target.value;
  renderPlanAndCandidates();
});

document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-position-sort]");
  if (!button) return;
  const key = button.dataset.positionSort;
  positionSort = {
    key,
    direction: positionSort.key === key && positionSort.direction === "desc" ? "asc" : positionSortDefaultDirection(key)
  };
  renderPositions(computePositions());
});

elements.positionMobileSort?.addEventListener("change", (event) => {
  positionSort = parseMobileSortValue(event.target.value);
  renderPositions(computePositions());
});

document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-sunny30-sort]");
  if (!button) return;
  const key = button.dataset.sunny30Sort;
  sunny30Sort = {
    key,
    direction: sunny30Sort.key === key
      ? (sunny30Sort.direction === "desc" ? "asc" : "desc")
      : sunny30DefaultDirection(key)
  };
  renderSunny30(computePositions());
});

elements.sunny30MobileSort?.addEventListener("change", (event) => {
  sunny30Sort = parseMobileSortValue(event.target.value);
  renderSunny30(computePositions());
});

document.addEventListener("click", (event) => {
  const positionButton = event.target.closest("[data-toggle-position-card]");
  if (positionButton) {
    const symbol = normalizeSymbol(positionButton.dataset.togglePositionCard);
    if (expandedPositionCards.has(symbol)) expandedPositionCards.delete(symbol);
    else expandedPositionCards.add(symbol);
    renderPositions(computePositions());
    return;
  }

  const sunny30Button = event.target.closest("[data-toggle-sunny30-card]");
  if (sunny30Button) {
    const symbol = normalizeSymbol(sunny30Button.dataset.toggleSunny30Card);
    if (expandedSunny30Cards.has(symbol)) expandedSunny30Cards.delete(symbol);
    else expandedSunny30Cards.add(symbol);
    renderSunny30(computePositions());
    return;
  }

});

document.addEventListener("change", async (event) => {
  const select = event.target.closest("[data-stock-category-select]");
  if (!select) return;
  const symbol = normalizeSymbol(select.dataset.symbol);
  const previous = select.dataset.currentCategory || select.value;
  try {
    select.disabled = true;
    const category = await saveStockCategory(symbol, select.value);
    setQuoteUpdateStatus(`已更新 ${symbol} 分类：${category}`, "success");
  } catch (error) {
    select.value = previous;
    select.disabled = false;
    window.alert(error.message);
  }
});

document.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-delete-sunny30]");
  if (!button) return;
  const symbol = normalizeSymbol(button.dataset.deleteSunny30);
  const candidate = (state.candidates ?? []).find((item) => normalizeSymbol(item.symbol) === symbol);
  const name = candidate?.name || symbol;
  if (!window.confirm(`从晴仓30删除「${name}」？这里只移出晴仓30，不会删除持仓。`)) return;

  try {
    button.disabled = true;
    await deleteSunny30Candidate(symbol);
    setQuoteUpdateStatus(`已从晴仓30删除 ${name}`, "success");
  } catch (error) {
    button.disabled = false;
    window.alert(error.message);
  }
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

elements.masterMatrixFilters?.addEventListener("click", (event) => {
  const button = event.target.closest("[data-master-matrix-filter]");
  if (!button) return;
  masterMatrixFilter = button.dataset.masterMatrixFilter || "all";
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

document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-detail-section]");
  if (!button) return;
  const sectionId = button.dataset.detailSection;
  expandStockDetailSection(sectionId);
  setStockDetailActiveSection(sectionId);
  document.getElementById(sectionId)?.scrollIntoView({ behavior: "smooth", block: "start" });
});

document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-stock-detail-toggle]");
  if (!button) return;
  toggleStockDetailSection(button.dataset.stockDetailToggle);
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

document.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-update-industries]");
  if (!button) return;
  try {
    await updateIndustryMetrics(button);
  } catch (error) {
    setQuoteUpdateStatus(error.message, "error");
  }
});

document.addEventListener("submit", async (event) => {
  const form = event.target.closest("[data-screening-weights-form]");
  if (!form) return;
  event.preventDefault();
  try {
    await saveScreeningWeights(new FormData(form));
    setQuoteUpdateStatus("选股排序权重已保存", "success");
  } catch (error) {
    window.alert(error.message);
  }
});

document.addEventListener("submit", async (event) => {
  const form = event.target.closest("[data-stock-human-input-form]");
  if (!form) return;
  event.preventDefault();
  try {
    await saveStockHumanInputs(form);
    setQuoteUpdateStatus("人工判断已保存", "success");
  } catch (error) {
    window.alert(error.message);
  }
});

document.addEventListener("submit", async (event) => {
  const form = event.target.closest("[data-stock-valuation-form]");
  if (!form) return;
  event.preventDefault();
  try {
    await saveStockValuationInputs(form);
    setQuoteUpdateStatus("估值假设已保存", "success");
  } catch (error) {
    window.alert(error.message);
  }
});

document.addEventListener("submit", async (event) => {
  const form = event.target.closest("[data-stock-hold-log-form]");
  if (!form) return;
  event.preventDefault();
  try {
    await saveStockHoldLog(form);
    setQuoteUpdateStatus("决策日志已保存", "success");
  } catch (error) {
    window.alert(error.message);
  }
});

document.addEventListener("click", async (event) => {
  const button = event.target.closest("[data-create-hold-log]");
  if (!button) return;
  const symbol = normalizeSymbol(button.dataset.createHoldLog);
  const stock = auditUniverse(computePositions()).find((item) => normalizeSymbol(item.symbol) === symbol);
  const reason = window.prompt(`记录「${stock?.name || symbol}」继续持有理由`);
  if (!reason || !reason.trim()) {
    window.alert("继续持有日志必须填写人工理由");
    return;
  }
  try {
    button.disabled = true;
    const range = stock ? valuationRangeView(stock) : {};
    const marginText = Number.isFinite(range.margin) ? `安全边际 ${percent(range.margin * 100, false)}` : "安全边际待补";
    setLoadedState(await requestJSON("/api/decision-logs", {
      method: "POST",
      body: JSON.stringify({
        type: "hold",
        symbol,
        decision: "继续持有",
        detail: `${reason.trim()}；${marginText}`
      })
    }));
    localStorage.removeItem(STORAGE_KEY);
    render();
    setQuoteUpdateStatus(`已记录 ${stock?.name || symbol} 继续持有理由`, "success");
  } catch (error) {
    window.alert(error.message);
  } finally {
    button.disabled = false;
  }
});

document.querySelector("#openTradePanelSecondary").addEventListener("click", () => {
  setTradeAssetType();
  elements.tradeDialog.showModal();
});

function openResearchDialog() {
  pendingResearch = null;
  elements.importResearchButton.disabled = true;
  elements.researchPreview.innerHTML = "";
  setResearchStatus("");
  elements.researchDialog.showModal();
}

document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-open-research]");
  if (!button) return;
  openResearchDialog();
});

document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-research-desk-filter]");
  if (!button) return;
  researchDeskFilter = button.dataset.researchDeskFilter || "all";
  renderIndustryDesk(computePositions());
});

document.addEventListener("click", (event) => {
  const button = event.target.closest("[data-research-desk-toggle]");
  if (!button) return;
  const section = button.dataset.researchDeskToggle;
  if (expandedResearchDeskSections.has(section)) {
    expandedResearchDeskSections.delete(section);
  } else {
    expandedResearchDeskSections.add(section);
  }
  renderIndustryDesk(computePositions());
});

elements.updateQuotesButton.addEventListener("click", async () => {
  try {
    await updateQuotes();
  } catch (error) {
    setQuoteUpdateStatus(error.message, "error");
  }
});

elements.updateValuationHistoryButton?.addEventListener("click", async () => {
  try {
    await updateValuationHistory();
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

elements.openSunny30CandidateButton?.addEventListener("click", openSunny30CandidateDialog);

document.querySelector("#closeSunny30CandidatePanel")?.addEventListener("click", () => {
  elements.sunny30CandidateDialog.close();
});

document.querySelector("#cancelSunny30Candidate")?.addEventListener("click", () => {
  elements.sunny30CandidateDialog.close();
});

document.querySelector("#closeResearchPanel").addEventListener("click", () => {
  elements.researchDialog.close();
});

document.querySelector("#cancelResearch").addEventListener("click", () => {
  elements.researchDialog.close();
});

elements.tradeForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    await addTrade(new FormData(elements.tradeForm));
    elements.tradeForm.reset();
    elements.tradeDialog.close();
  } catch (error) {
    window.alert(error.message);
  }
});

elements.holdingForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  await saveHolding(new FormData(elements.holdingForm));
  elements.holdingDialog.close();
});

elements.sunny30CandidateForm?.addEventListener("submit", async (event) => {
  event.preventDefault();
  try {
    await saveSunny30Candidate(new FormData(elements.sunny30CandidateForm));
    elements.sunny30CandidateForm.reset();
    elements.sunny30CandidateDialog.close();
    setQuoteUpdateStatus("晴仓30跟踪标的已保存", "success");
  } catch (error) {
    window.alert(error.message);
  }
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
  setPrivacyToggleState();
  if (USE_BACKEND) {
    renderLoadingState();
  }
  const loaded = await loadBackendState();
  if (!loaded && USE_BACKEND) {
    setQuoteUpdateStatus("后端不可用，使用浏览器本地兜底数据", "error");
  }
  syncCash();
  handleRoute(window.location.hash.slice(1));
}

init();
