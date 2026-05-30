package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"gorm.io/gorm"
)

type categorySeed struct {
	name        string
	color       string
	description string
	children    []categorySeed
}

var categorySeeds = []categorySeed{
	{name: "工作", color: "#EF4444", description: "工作相关的时间记录", children: []categorySeed{
		{name: "项目管理", color: "#EF4444", description: "项目规划与跟进"},
	}},
	{name: "学习", color: "#10B981", description: "学习和培训时间"},
	{name: "会议", color: "#F59E0B", description: "各种会议时间"},
	{name: "开发", color: "#8B5CF6", description: "软件开发和编程", children: []categorySeed{
		{name: "前端", color: "#8B5CF6", description: "前端开发"},
		{name: "后端", color: "#7C3AED", description: "后端开发"},
	}},
	{name: "休息", color: "#6B7280", description: "休息和放松时间"},
	{name: "运动", color: "#F97316", description: "体育锻炼和健身"},
	{name: "其他", color: "#6366F1", description: "其他未分类活动"},
}

var remarksByCategory = map[string][]string{
	"工作": {
		"处理客户反馈的紧急需求",
		"编写项目周报并同步给团队",
		"Review 同事提交的 PR",
		"整理 Q2 季度工作计划",
		"与产品经理对接需求文档",
		"优化内部工作流程文档",
		"处理邮件和日常沟通",
		"更新项目进度看板",
	},
	"项目管理": {
		"更新项目里程碑",
		"整理风险清单",
		"同步跨团队依赖",
	},
	"学习": {
		"阅读《设计模式》第 3 章",
		"学习 Go 语言并发编程",
		"完成在线课程第 5 节作业",
		"研究 Vue 3 新特性 Composition API",
		"复习数据库索引优化相关知识",
		"学习 Kubernetes 基础概念",
		"阅读技术博客并做笔记",
		"练习算法题 3 道",
	},
	"会议": {
		"每日站会同步开发进度",
		"技术方案评审会议",
		"与后端团队对接 API 接口",
		"项目复盘会议",
		"季度 OKR 对齐会",
		"客户需求沟通会",
		"团队代码规范讨论",
		"跨部门协作沟通",
	},
	"开发": {
		"开发用户登录模块",
		"调试支付接口回调逻辑",
		"重构日志模块代码",
		"编写单元测试覆盖核心逻辑",
		"优化数据库查询性能",
		"实现文件上传功能",
		"修复生产环境 Bug #428",
		"集成第三方 OAuth 登录",
	},
	"前端": {
		"实现任务列表筛选组件",
		"调整分类树样式",
		"修复表单校验提示",
	},
	"后端": {
		"编写 mapper 单元测试",
		"调整 API 响应结构",
		"优化 timelog 查询 SQL",
	},
	"休息": {
		"午休小憩恢复精力",
		"下午茶时间放松身心",
		"冥想 15 分钟缓解压力",
		"散步休息眼睛",
		"听音乐放松",
		"午休后泡杯咖啡",
		"和同事闲聊放松",
		"闭目养神",
	},
	"运动": {
		"晨跑 5 公里",
		"健身房力量训练",
		"瑜伽拉伸 30 分钟",
		"傍晚散步 40 分钟",
		"跳绳 20 分钟",
		"骑自行车通勤",
		"篮球训练",
		"游泳 1000 米",
	},
	"其他": {
		"整理书桌和文件",
		"去超市采购生活用品",
		"做饭准备晚餐",
		"看一集纪录片",
		"打扫房间卫生",
		"给植物浇水",
		"规划周末行程",
		"记账和财务整理",
	},
}

type taskSeed struct {
	title            string
	description      string
	category         string
	dueOffsetDays    int
	estimatedMinutes int32
	completed        bool
	suspended        bool
}

var taskSeeds = []taskSeed{
	{title: "整理本周工作周报", category: "工作", dueOffsetDays: -2, estimatedMinutes: 45, completed: true},
	{title: "跟进客户 A 需求变更", category: "工作", dueOffsetDays: 1, estimatedMinutes: 90},
	{title: "更新项目风险清单", category: "项目管理", dueOffsetDays: 0, estimatedMinutes: 60},
	{title: "完成 Go 并发章节笔记", category: "学习", dueOffsetDays: -1, estimatedMinutes: 120, completed: true},
	{title: "刷 5 道 LeetCode 中等题", category: "学习", dueOffsetDays: 2, estimatedMinutes: 90},
	{title: "准备技术方案评审材料", category: "会议", dueOffsetDays: 0, estimatedMinutes: 60},
	{title: "组织季度 OKR 复盘", category: "会议", dueOffsetDays: 5, estimatedMinutes: 120, suspended: true},
	{title: "实现 timelog 表单校验", category: "开发", dueOffsetDays: -3, estimatedMinutes: 180, completed: true},
	{title: "补充 contract test 用例", category: "开发", dueOffsetDays: 1, estimatedMinutes: 120},
	{title: "优化分类树拖拽交互", category: "前端", dueOffsetDays: 3, estimatedMinutes: 150},
	{title: "调整 task 筛选 API", category: "后端", dueOffsetDays: 0, estimatedMinutes: 90},
	{title: "修复 passkey 注册 404", category: "后端", dueOffsetDays: -1, estimatedMinutes: 60, completed: true},
	{title: "安排一次长散步", category: "休息", dueOffsetDays: 1, estimatedMinutes: 40},
	{title: "周末晨跑计划", category: "运动", dueOffsetDays: 4, estimatedMinutes: 60},
	{title: "整理家庭账单", category: "其他", dueOffsetDays: 6, estimatedMinutes: 45, suspended: true},
}

var durations = []time.Duration{
	30 * time.Minute,
	45 * time.Minute,
	60 * time.Minute,
	90 * time.Minute,
	2 * time.Hour,
	3 * time.Hour,
}

func main() {
	config.ResetConfig()
	cfg := config.GetConfig("./config.yml")

	log := logger.SetZapLogger(*cfg)

	dao, err := model.NewDao(cfg, log)
	if err != nil {
		panic(err)
	}
	db := dao.Db()

	if err := resetSeedData(db); err != nil {
		panic(err)
	}

	categoryIDs, err := seedCategories(db)
	if err != nil {
		panic(err)
	}

	taskCount, tasksByCategory, err := seedTasks(db, categoryIDs)
	if err != nil {
		panic(err)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc)
	startDate := now.AddDate(0, 0, -28)

	r := rand.New(rand.NewSource(42))
	timelogCount, linkedCount := seedTimelogs(db, categoryIDs, tasksByCategory, startDate, loc, r)

	fmt.Printf("Seeded %d categories, %d tasks, %d timelogs (%d linked to tasks)\n",
		len(categoryIDs), taskCount, timelogCount, linkedCount)
}

func resetSeedData(db *gorm.DB) error {
	if err := db.Exec("DELETE FROM timelogs").Error; err != nil {
		return fmt.Errorf("clear timelogs: %w", err)
	}
	if err := db.Exec("DELETE FROM tasks").Error; err != nil {
		return fmt.Errorf("clear tasks: %w", err)
	}
	if err := db.Exec("DELETE FROM categories").Error; err != nil {
		return fmt.Errorf("clear categories: %w", err)
	}
	return nil
}

func seedCategories(db *gorm.DB) (map[string]int32, error) {
	ids := make(map[string]int32)
	for _, seed := range categorySeeds {
		if err := createCategoryTree(db, seed, nil, ids); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

func createCategoryTree(db *gorm.DB, seed categorySeed, parentID *int32, ids map[string]int32) error {
	category := &gen.Category{
		Name:        seed.name,
		Color:       strPtr(seed.color),
		Description: strPtr(seed.description),
	}
	if parentID != nil {
		category.ParentID = parentID
	}
	if err := model.CreateCategory(db, category); err != nil {
		return fmt.Errorf("create category %q: %w", seed.name, err)
	}
	ids[seed.name] = *category.ID
	for _, child := range seed.children {
		if err := createCategoryTree(db, child, category.ID, ids); err != nil {
			return err
		}
	}
	return nil
}

func seedTasks(db *gorm.DB, categoryIDs map[string]int32) (int, map[int32][]*gen.Task, error) {
	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc)
	tasksByCategory := make(map[int32][]*gen.Task)

	for _, seed := range taskSeeds {
		categoryID, ok := categoryIDs[seed.category]
		if !ok {
			return 0, nil, fmt.Errorf("unknown task category %q", seed.category)
		}

		dueDate := dateAt(now.AddDate(0, 0, seed.dueOffsetDays), loc)
		task := &gen.Task{
			Title:            seed.title,
			Description:      strPtr(seed.description),
			CategoryID:       categoryID,
			DueDate:          dueDate,
			EstimatedMinutes: seed.estimatedMinutes,
			IsCompleted:      boolPtr(seed.completed),
			IsSuspended:      boolPtr(seed.suspended),
		}
		if seed.completed {
			completedAt := dueDate.Add(2 * time.Hour)
			task.CompletedAt = &completedAt
		}

		if err := model.CreateTask(db, task); err != nil {
			return 0, nil, fmt.Errorf("create task %q: %w", seed.title, err)
		}
		tasksByCategory[categoryID] = append(tasksByCategory[categoryID], task)
	}

	return len(taskSeeds), tasksByCategory, nil
}

func seedTimelogs(
	db *gorm.DB,
	categoryIDs map[string]int32,
	tasksByCategory map[int32][]*gen.Task,
	startDate time.Time,
	loc *time.Location,
	r *rand.Rand,
) (int, int) {
	workdayCategories := []int32{
		categoryIDs["工作"],
		categoryIDs["学习"],
		categoryIDs["会议"],
		categoryIDs["开发"],
	}
	weekendCategories := []int32{
		categoryIDs["休息"],
		categoryIDs["运动"],
		categoryIDs["其他"],
	}
	allCategories := append(append([]int32{}, workdayCategories...), weekendCategories...)
	allCategories = append(allCategories, categoryIDs["前端"], categoryIDs["后端"], categoryIDs["项目管理"])

	var total, linked int
	for day := 0; day <= 28; day++ {
		currentDay := startDate.AddDate(0, 0, day)
		if r.Float32() < 0.15 {
			continue
		}

		hourOffset := 8 + r.Intn(3)
		minuteOffset := r.Intn(60)
		currentTime := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), hourOffset, minuteOffset, 0, 0, loc)

		numLogs := 3 + r.Intn(5)
		for i := 0; i < numLogs; i++ {
			catID := allCategories[r.Intn(len(allCategories))]
			weekday := currentDay.Weekday()
			if weekday >= time.Monday && weekday <= time.Friday {
				if r.Float32() < 0.6 {
					catID = workdayCategories[r.Intn(len(workdayCategories))]
				}
			} else if r.Float32() < 0.5 {
				catID = weekendCategories[r.Intn(len(weekendCategories))]
			}

			dur := durations[r.Intn(len(durations))]
			start := currentTime
			end := start.Add(dur)
			if end.Day() != start.Day() {
				end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 0, 0, loc)
			}

			remark := pickRemark(catID, categoryIDs, r)
			tl := &gen.Timelog{
				UserID:     int32Ptr(1),
				StartTime:  start.UTC(),
				EndTime:    timePtr(end.UTC()),
				CategoryID: catID,
				Remark:     &remark,
			}

			if r.Float32() < 0.25 {
				if task := pickTask(tasksByCategory[catID], r); task != nil {
					tl.TaskID = task.ID
					linked++
				}
			}

			if err := db.Create(tl).Error; err != nil {
				fmt.Printf("Error creating timelog: %v\n", err)
				continue
			}
			total++

			currentTime = end.Add(time.Duration(r.Intn(30)) * time.Minute)
			if currentTime.Hour() >= 23 {
				break
			}
		}
	}

	return total, linked
}

func pickRemark(categoryID int32, categoryIDs map[string]int32, r *rand.Rand) string {
	for name, id := range categoryIDs {
		if id != categoryID {
			continue
		}
		remarks := remarksByCategory[name]
		if len(remarks) == 0 {
			return name
		}
		return remarks[r.Intn(len(remarks))]
	}
	return "时间记录"
}

func pickTask(tasks []*gen.Task, r *rand.Rand) *gen.Task {
	if len(tasks) == 0 {
		return nil
	}
	return tasks[r.Intn(len(tasks))]
}

func dateAt(day time.Time, loc *time.Location) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), 9, 0, 0, 0, loc).UTC()
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(v bool) *bool {
	return &v
}

func int32Ptr(i int32) *int32 {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}
