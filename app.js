const STORAGE_KEY = "stock-portfolio-desk-v2";

const seedState = {
  totalCapital: 1100000,
  cash: 0,
  fx: { CNY: 1, HKD: 0.8716, USD: 7.1 },
  trades: [],
  holdings: [
    {
      symbol: "0700.HK",
      name: "腾讯控股",
      shares: 200,
      cost: 480.43,
      currentPrice: 463,
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
    { symbol: "600036.SH", name: "招商银行", shares: 500, cost: 39.18, currentPrice: 39.18, currency: "CNY", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "0696.HK", name: "民航信", shares: 11000, cost: 10.648, currentPrice: 10.648, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "0506.HK", name: "中国食品", shares: 22000, cost: 4.041, currentPrice: 4.041, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "2669.HK", name: "中海物业", shares: 20000, cost: 4.468, currentPrice: 4.468, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "6049.HK", name: "保利物业", shares: 2600, cost: 32.663, currentPrice: 32.663, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" },
    { symbol: "0883.HK", name: "中海油", shares: 2000, cost: 29.326, currentPrice: 29.326, currency: "HKD", action: "", status: "", marginOfSafety: null, qualityScore: null, industry: "", notes: "" }
  ],
  plan: [
    { rank: 1, name: "腾讯控股", priority: "观察/低优先级", advice: "继续持有；新资金等待≤HK$432，HK$400-430可分批", discipline: "优秀资产要求≥15%安全边际；当前约9%，未达标" },
    { rank: 2, name: "美的集团", priority: "核心替补/中优先级", advice: "A股等待≤¥76分批；H股≤HK$86-87优先；当前不追买", discipline: "优秀资产要求≥20%安全边际；A股当前约15.3%，未达标" },
    { rank: 3, name: "海康威视", priority: "重点预期差候选/中优先级", advice: "不重仓；¥35-37仅适合小仓验证，¥30-32更从容；Q2验证后可升核心替补", discipline: "质量分84，合格候选要求≥25%安全边际" },
    { rank: 4, name: "伊利股份", priority: "核心替补/中低优先级", advice: "暂不追买；¥25-26开始关注，≤¥24可考虑分批", discipline: "质量分83，合格候选要求≥25%安全边际" }
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

let state = loadState();
let activeFilter = "all";
let searchTerm = "";
const pageTitles = {
  overview: "股票持仓管理",
  positions: "当前持仓",
  trades: "交易流水",
  research: "研究与调仓计划",
  candidates: "候选池"
};

const elements = {
  pageTitle: document.querySelector("#pageTitle"),
  positionsBody: document.querySelector("#positionsBody"),
  tradeList: document.querySelector("#tradeList"),
  allocationChart: document.querySelector("#allocationChart"),
  overviewPlanList: document.querySelector("#overviewPlanList"),
  planList: document.querySelector("#planList"),
  candidateList: document.querySelector("#candidateList"),
  rulesList: document.querySelector("#rulesList"),
  totalValue: document.querySelector("#totalValue"),
  totalPnl: document.querySelector("#totalPnl"),
  totalPnlRate: document.querySelector("#totalPnlRate"),
  dayChange: document.querySelector("#dayChange"),
  dayChangeRate: document.querySelector("#dayChangeRate"),
  cashBalance: document.querySelector("#cashBalance"),
  positionCount: document.querySelector("#positionCount"),
  recordCount: document.querySelector("#recordCount"),
  searchInput: document.querySelector("#searchInput"),
  tradeDialog: document.querySelector("#tradeDialog"),
  tradeForm: document.querySelector("#tradeForm"),
  holdingDialog: document.querySelector("#holdingDialog"),
  holdingForm: document.querySelector("#holdingForm")
};

syncCash();

function loadState() {
  const saved = localStorage.getItem(STORAGE_KEY);
  if (!saved) return structuredClone(seedState);

  try {
    return { ...structuredClone(seedState), ...JSON.parse(saved) };
  } catch {
    return structuredClone(seedState);
  }
}

function saveState() {
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
    saveState();
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

function computePositions() {
  return state.holdings
    .filter((holding) => holding.shares > 0)
    .map((holding) => {
      const marketValueLocal = holding.shares * holding.currentPrice;
      const costValueLocal = holding.shares * holding.cost;
      const marketValueCny = marketValueLocal * fx(holding.currency);
      const costValueCny = costValueLocal * fx(holding.currency);
      const pnlCny = marketValueCny - costValueCny;

      return {
        ...holding,
        marketValueLocal,
        marketValueCny,
        costValueCny,
        pnlCny,
        pnlRate: costValueCny ? (pnlCny / costValueCny) * 100 : 0,
        dayChange: marketValueCny * ((Math.sin(holding.symbol.charCodeAt(0)) * 0.018) + 0.004)
      };
    });
}

function syncCash() {
  const invested = state.holdings.reduce((sum, holding) => {
    return sum + holding.shares * holding.currentPrice * fx(holding.currency);
  }, 0);

  if (!state.trades.length) {
    state.cash = state.totalCapital - invested;
  }
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
    elements.positionsBody.innerHTML = `<tr><td colspan="10" class="empty-state">暂无符合条件的持仓</td></tr>`;
    return;
  }

  elements.positionsBody.innerHTML = filtered
    .map((position) => {
      const weight = totalValue ? (position.marketValueCny / totalValue) * 100 : 0;
      const pnlClass = position.pnlCny >= 0 ? "positive" : "negative";
      const marginText = Number.isFinite(position.marginOfSafety)
        ? percent(position.marginOfSafety * 100, false)
        : "-";
      const qualityText = Number.isFinite(position.qualityScore) ? position.qualityScore : "-";

      return `
        <tr>
          <td>
            <div class="stock-cell">
              <span class="ticker">${position.symbol.slice(0, 4)}</span>
              <span class="stock-name">
                <strong>${position.name}</strong>
                <span>${position.symbol} · ${position.industry || "未分类"}</span>
              </span>
            </div>
          </td>
          <td>${position.shares}</td>
          <td>${currency(position.cost, position.currency)}</td>
          <td>${currency(position.currentPrice, position.currency)}</td>
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
            <button class="icon-button edit-holding" data-symbol="${position.symbol}" title="编辑 Excel 信息">✎</button>
          </td>
        </tr>
        <tr class="detail-row">
          <td colspan="10">
            <strong>${position.status || "未填写达标状态"}</strong>
            <span>${position.action || "未填写最终动作"}</span>
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
        <div class="allocation-row">
          <strong>${position.symbol}</strong>
          <div class="allocation-track">
            <span style="width: ${share}%; background: ${palette[index % palette.length]}"></span>
          </div>
          <span>${share.toFixed(1)}%</span>
        </div>
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

function renderResearch() {
  elements.overviewPlanList.innerHTML = state.plan
    .slice(0, 4)
    .map((item) => `
      <div class="research-item">
        <strong>${item.rank}. ${item.name}</strong>
        <span>${item.priority}</span>
        <small>${item.advice}</small>
      </div>
    `)
    .join("");

  elements.planList.innerHTML = state.plan
    .map((item) => `
      <div class="research-item">
        <strong>${item.rank}. ${item.name}</strong>
        <span>${item.priority}</span>
        <small>${item.advice}</small>
      </div>
    `)
    .join("");

  elements.candidateList.innerHTML = state.candidates
    .map((item) => `
      <div class="research-item">
        <strong>${item.name} · ${item.symbol}</strong>
        <span>${item.action}</span>
        <small>质量 ${item.qualityScore} · 安全边际 ${percent(item.marginOfSafety * 100, false)} · 买入价 ${currency(item.targetBuyPrice, item.currency)}</small>
      </div>
    `)
    .join("");

  elements.rulesList.innerHTML = state.rules
    .map((item) => `
      <div class="rule-item">
        <strong>${item.dimension}</strong>
        <span>${item.score} 分</span>
        <small>${item.standard}</small>
      </div>
    `)
    .join("");
}

function render() {
  const positions = computePositions();
  renderMetrics(positions);
  renderPositions(positions);
  renderAllocation(positions);
  renderTrades();
  renderResearch();
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
  state.trades.push(trade);
  state.cash += side === "sell" ? shares * price * fx(currencyCode) : -(shares * price * fx(currencyCode));
  saveState();
  render();
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
  form.marginOfSafety.value = Number.isFinite(holding.marginOfSafety) ? holding.marginOfSafety : "";
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
    marginOfSafety: formData.get("marginOfSafety") === "" ? null : Number(formData.get("marginOfSafety")),
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
  const nextView = pageTitles[view] ? view : "overview";
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

document.querySelector("#openTradePanel").addEventListener("click", () => {
  elements.tradeDialog.showModal();
});

document.querySelector("#openTradePanelSecondary").addEventListener("click", () => {
  elements.tradeDialog.showModal();
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

document.querySelector("#seedButton").addEventListener("click", async () => {
  if (USE_BACKEND) {
    state = await requestJSON("/api/reset", { method: "POST" });
  } else {
    state = structuredClone(seedState);
    syncCash();
  }
  saveState();
  render();
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
