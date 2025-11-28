package builder

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/xiaobei/singbox-manager/internal/storage"
)

// SingBoxConfig sing-box 配置结构
type SingBoxConfig struct {
	Log          *LogConfig          `json:"log,omitempty"`
	DNS          *DNSConfig          `json:"dns,omitempty"`
	NTP          *NTPConfig          `json:"ntp,omitempty"`
	Inbounds     []Inbound           `json:"inbounds,omitempty"`
	Outbounds    []Outbound          `json:"outbounds"`
	Route        *RouteConfig        `json:"route,omitempty"`
	Experimental *ExperimentalConfig `json:"experimental,omitempty"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level     string `json:"level,omitempty"`
	Timestamp bool   `json:"timestamp,omitempty"`
	Output    string `json:"output,omitempty"`
}

// DNSConfig DNS 配置
type DNSConfig struct {
	Servers []DNSServer `json:"servers,omitempty"`
	Rules   []DNSRule   `json:"rules,omitempty"`
	Final   string      `json:"final,omitempty"`
}

// DNSServer DNS 服务器
type DNSServer struct {
	Tag            string `json:"tag"`
	Address        string `json:"address"`
	AddressResolver string `json:"address_resolver,omitempty"`
	Detour         string `json:"detour,omitempty"`
}

// DNSRule DNS 规则
type DNSRule struct {
	RuleSet   []string `json:"rule_set,omitempty"`
	QueryType []string `json:"query_type,omitempty"`
	Server    string   `json:"server,omitempty"`
	Outbound  string   `json:"outbound,omitempty"`
}

// NTPConfig NTP 配置
type NTPConfig struct {
	Enabled bool   `json:"enabled"`
	Server  string `json:"server,omitempty"`
}

// Inbound 入站配置
type Inbound struct {
	Type           string   `json:"type"`
	Tag            string   `json:"tag"`
	Listen         string   `json:"listen,omitempty"`
	ListenPort     int      `json:"listen_port,omitempty"`
	Address        []string `json:"address,omitempty"`
	AutoRoute      bool     `json:"auto_route,omitempty"`
	StrictRoute    bool     `json:"strict_route,omitempty"`
	Stack          string   `json:"stack,omitempty"`
	Sniff          bool     `json:"sniff,omitempty"`
	SniffOverrideDestination bool `json:"sniff_override_destination,omitempty"`
}

// Outbound 出站配置
type Outbound map[string]interface{}

// RouteConfig 路由配置
type RouteConfig struct {
	Rules               []RouteRule `json:"rules,omitempty"`
	RuleSet             []RuleSet   `json:"rule_set,omitempty"`
	Final               string      `json:"final,omitempty"`
	AutoDetectInterface bool        `json:"auto_detect_interface,omitempty"`
}

// RouteRule 路由规则
type RouteRule map[string]interface{}

// RuleSet 规则集
type RuleSet struct {
	Tag            string `json:"tag"`
	Type           string `json:"type"`
	Format         string `json:"format"`
	URL            string `json:"url,omitempty"`
	DownloadDetour string `json:"download_detour,omitempty"`
}

// ExperimentalConfig 实验性配置
type ExperimentalConfig struct {
	ClashAPI *ClashAPIConfig `json:"clash_api,omitempty"`
	CacheFile *CacheFileConfig `json:"cache_file,omitempty"`
}

// ClashAPIConfig Clash API 配置
type ClashAPIConfig struct {
	ExternalController string `json:"external_controller,omitempty"`
	ExternalUI         string `json:"external_ui,omitempty"`
	ExternalUIDownloadURL string `json:"external_ui_download_url,omitempty"`
	Secret             string `json:"secret,omitempty"`
	DefaultMode        string `json:"default_mode,omitempty"`
}

// CacheFileConfig 缓存文件配置
type CacheFileConfig struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path,omitempty"`
}

// ConfigBuilder 配置生成器
type ConfigBuilder struct {
	settings   *storage.Settings
	nodes      []storage.Node
	filters    []storage.Filter
	rules      []storage.Rule
	ruleGroups []storage.RuleGroup
}

// NewConfigBuilder 创建配置生成器
func NewConfigBuilder(settings *storage.Settings, nodes []storage.Node, filters []storage.Filter, rules []storage.Rule, ruleGroups []storage.RuleGroup) *ConfigBuilder {
	return &ConfigBuilder{
		settings:   settings,
		nodes:      nodes,
		filters:    filters,
		rules:      rules,
		ruleGroups: ruleGroups,
	}
}

// Build 构建 sing-box 配置
func (b *ConfigBuilder) Build() (*SingBoxConfig, error) {
	config := &SingBoxConfig{
		Log:       b.buildLog(),
		DNS:       b.buildDNS(),
		NTP:       b.buildNTP(),
		Inbounds:  b.buildInbounds(),
		Outbounds: b.buildOutbounds(),
		Route:     b.buildRoute(),
	}

	// 添加 Clash API 支持
	if b.settings.ClashAPIPort > 0 {
		config.Experimental = b.buildExperimental()
	}

	return config, nil
}

// BuildJSON 构建 JSON 字符串
func (b *ConfigBuilder) BuildJSON() (string, error) {
	config, err := b.Build()
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化配置失败: %w", err)
	}

	return string(data), nil
}

// buildLog 构建日志配置
func (b *ConfigBuilder) buildLog() *LogConfig {
	return &LogConfig{
		Level:     "info",
		Timestamp: true,
	}
}

// buildDNS 构建 DNS 配置
func (b *ConfigBuilder) buildDNS() *DNSConfig {
	return &DNSConfig{
		Servers: []DNSServer{
			{
				Tag:            "dns_proxy",
				Address:        b.settings.ProxyDNS,
				AddressResolver: "dns_resolver",
				Detour:         "Proxy",
			},
			{
				Tag:            "dns_direct",
				Address:        b.settings.DirectDNS,
				AddressResolver: "dns_resolver",
				Detour:         "DIRECT",
			},
			{
				Tag:     "dns_resolver",
				Address: "223.5.5.5",
			},
			{
				Tag:     "dns_block",
				Address: "rcode://success",
			},
		},
		Rules: []DNSRule{
			{
				RuleSet: []string{"geosite-category-ads-all"},
				Server:  "dns_block",
			},
			{
				RuleSet: []string{"geosite-geolocation-cn"},
				Server:  "dns_direct",
			},
			{
				RuleSet: []string{"geosite-geolocation-!cn"},
				Server:  "dns_proxy",
			},
		},
		Final: "dns_direct",
	}
}

// buildNTP 构建 NTP 配置
func (b *ConfigBuilder) buildNTP() *NTPConfig {
	return &NTPConfig{
		Enabled: true,
		Server:  "time.apple.com",
	}
}

// buildInbounds 构建入站配置
func (b *ConfigBuilder) buildInbounds() []Inbound {
	inbounds := []Inbound{
		{
			Type:       "mixed",
			Tag:        "mixed-in",
			Listen:     "127.0.0.1",
			ListenPort: b.settings.MixedPort,
			Sniff:      true,
			SniffOverrideDestination: true,
		},
	}

	if b.settings.TunEnabled {
		inbounds = append(inbounds, Inbound{
			Type:        "tun",
			Tag:         "tun-in",
			Address:     []string{"172.19.0.1/30", "fdfe:dcba:9876::1/126"},
			AutoRoute:   true,
			StrictRoute: true,
			Stack:       "system",
			Sniff:       true,
			SniffOverrideDestination: true,
		})
	}

	return inbounds
}

// buildOutbounds 构建出站配置
func (b *ConfigBuilder) buildOutbounds() []Outbound {
	outbounds := []Outbound{
		{"type": "direct", "tag": "DIRECT"},
		{"type": "block", "tag": "REJECT"},
		{"type": "dns", "tag": "dns-out"},
	}

	// 收集所有节点标签和按国家分组
	var allNodeTags []string
	nodeTagSet := make(map[string]bool)
	countryNodes := make(map[string][]string) // 国家代码 -> 节点标签列表

	// 添加所有节点
	for _, node := range b.nodes {
		outbound := b.nodeToOutbound(node)
		outbounds = append(outbounds, outbound)
		tag := node.Tag
		if !nodeTagSet[tag] {
			allNodeTags = append(allNodeTags, tag)
			nodeTagSet[tag] = true
		}

		// 按国家分组
		if node.Country != "" {
			countryNodes[node.Country] = append(countryNodes[node.Country], tag)
		}
	}

	// 收集过滤器分组
	var filterGroupTags []string
	filterNodeMap := make(map[string][]string)

	for _, filter := range b.filters {
		if !filter.Enabled {
			continue
		}

		// 根据过滤器筛选节点
		var filteredTags []string
		for _, node := range b.nodes {
			if b.matchFilter(node, filter) {
				filteredTags = append(filteredTags, node.Tag)
			}
		}

		if len(filteredTags) == 0 {
			continue
		}

		groupTag := filter.Name
		filterGroupTags = append(filterGroupTags, groupTag)
		filterNodeMap[groupTag] = filteredTags

		// 创建分组
		group := Outbound{
			"tag":       groupTag,
			"type":      filter.Mode,
			"outbounds": filteredTags,
		}

		if filter.Mode == "urltest" {
			if filter.URLTestConfig != nil {
				group["url"] = filter.URLTestConfig.URL
				group["interval"] = filter.URLTestConfig.Interval
				group["tolerance"] = filter.URLTestConfig.Tolerance
			} else {
				group["url"] = "https://www.gstatic.com/generate_204"
				group["interval"] = "5m"
				group["tolerance"] = 50
			}
		}

		outbounds = append(outbounds, group)
	}

	// 创建按国家分组的出站选择器
	var countryGroupTags []string
	// 按国家代码排序，确保顺序一致
	var countryCodes []string
	for code := range countryNodes {
		countryCodes = append(countryCodes, code)
	}
	sort.Strings(countryCodes)

	for _, code := range countryCodes {
		nodes := countryNodes[code]
		if len(nodes) == 0 {
			continue
		}

		// 创建国家分组标签，格式: "🇭🇰 香港" 或 "HK"
		emoji := storage.GetCountryEmoji(code)
		name := storage.GetCountryName(code)
		groupTag := fmt.Sprintf("%s %s", emoji, name)
		countryGroupTags = append(countryGroupTags, groupTag)

		// 创建自动选择分组
		outbounds = append(outbounds, Outbound{
			"tag":       groupTag,
			"type":      "urltest",
			"outbounds": nodes,
			"url":       "https://www.gstatic.com/generate_204",
			"interval":  "5m",
			"tolerance": 50,
		})
	}

	// 创建自动选择组（所有节点）
	if len(allNodeTags) > 0 {
		outbounds = append(outbounds, Outbound{
			"tag":       "Auto",
			"type":      "urltest",
			"outbounds": allNodeTags,
			"url":       "https://www.gstatic.com/generate_204",
			"interval":  "5m",
			"tolerance": 50,
		})
	}

	// 创建主选择器
	proxyOutbounds := []string{"Auto"}
	proxyOutbounds = append(proxyOutbounds, countryGroupTags...) // 添加国家分组
	proxyOutbounds = append(proxyOutbounds, filterGroupTags...)
	proxyOutbounds = append(proxyOutbounds, allNodeTags...)

	outbounds = append(outbounds, Outbound{
		"tag":       "Proxy",
		"type":      "selector",
		"outbounds": proxyOutbounds,
		"default":   "Auto",
	})

	// 为启用的规则组创建选择器
	for _, rg := range b.ruleGroups {
		if !rg.Enabled {
			continue
		}

		selectorOutbounds := []string{"Proxy", "Auto", "DIRECT", "REJECT"}
		selectorOutbounds = append(selectorOutbounds, countryGroupTags...) // 添加国家分组
		selectorOutbounds = append(selectorOutbounds, filterGroupTags...)
		selectorOutbounds = append(selectorOutbounds, allNodeTags...)

		outbounds = append(outbounds, Outbound{
			"tag":       rg.Name,
			"type":      "selector",
			"outbounds": selectorOutbounds,
			"default":   rg.Outbound,
		})
	}

	// 创建漏网规则选择器
	fallbackOutbounds := []string{"Proxy", "DIRECT"}
	fallbackOutbounds = append(fallbackOutbounds, countryGroupTags...) // 添加国家分组
	fallbackOutbounds = append(fallbackOutbounds, filterGroupTags...)
	outbounds = append(outbounds, Outbound{
		"tag":       "Final",
		"type":      "selector",
		"outbounds": fallbackOutbounds,
		"default":   b.settings.FinalOutbound,
	})

	return outbounds
}

// nodeToOutbound 将节点转换为出站配置
func (b *ConfigBuilder) nodeToOutbound(node storage.Node) Outbound {
	outbound := Outbound{
		"tag":         node.Tag,
		"type":        node.Type,
		"server":      node.Server,
		"server_port": node.ServerPort,
	}

	// 复制 Extra 字段
	for k, v := range node.Extra {
		outbound[k] = v
	}

	return outbound
}

// matchFilter 检查节点是否匹配过滤器
func (b *ConfigBuilder) matchFilter(node storage.Node, filter storage.Filter) bool {
	name := strings.ToLower(node.Tag)

	// 1. 检查国家包含条件
	if len(filter.IncludeCountries) > 0 {
		matched := false
		for _, country := range filter.IncludeCountries {
			if strings.EqualFold(node.Country, country) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 2. 检查国家排除条件
	for _, country := range filter.ExcludeCountries {
		if strings.EqualFold(node.Country, country) {
			return false
		}
	}

	// 3. 检查关键字包含条件
	if len(filter.Include) > 0 {
		matched := false
		for _, keyword := range filter.Include {
			if strings.Contains(name, strings.ToLower(keyword)) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 4. 检查关键字排除条件
	for _, keyword := range filter.Exclude {
		if strings.Contains(name, strings.ToLower(keyword)) {
			return false
		}
	}

	return true
}

// buildRoute 构建路由配置
func (b *ConfigBuilder) buildRoute() *RouteConfig {
	route := &RouteConfig{
		AutoDetectInterface: true,
		Final:               "Final",
	}

	// 构建规则集
	ruleSetMap := make(map[string]bool)
	var ruleSets []RuleSet

	// 从规则组收集需要的规则集
	for _, rg := range b.ruleGroups {
		if !rg.Enabled {
			continue
		}
		for _, sr := range rg.SiteRules {
			tag := fmt.Sprintf("geosite-%s", sr)
			if !ruleSetMap[tag] {
				ruleSetMap[tag] = true
				ruleSets = append(ruleSets, RuleSet{
					Tag:            tag,
					Type:           "remote",
					Format:         "binary",
					URL:            fmt.Sprintf("%s/geosite-%s.srs", b.settings.RuleSetBaseURL, sr),
					DownloadDetour: "DIRECT",
				})
			}
		}
		for _, ir := range rg.IPRules {
			tag := fmt.Sprintf("geoip-%s", ir)
			if !ruleSetMap[tag] {
				ruleSetMap[tag] = true
				ruleSets = append(ruleSets, RuleSet{
					Tag:            tag,
					Type:           "remote",
					Format:         "binary",
					URL:            fmt.Sprintf("%s/../rule-set-geoip/geoip-%s.srs", b.settings.RuleSetBaseURL, ir),
					DownloadDetour: "DIRECT",
				})
			}
		}
	}

	route.RuleSet = ruleSets

	// 构建路由规则
	var rules []RouteRule

	// DNS 劫持规则
	rules = append(rules, RouteRule{
		"protocol": "dns",
		"outbound": "dns-out",
	})

	// 按优先级排序自定义规则
	sortedRules := make([]storage.Rule, len(b.rules))
	copy(sortedRules, b.rules)
	sort.Slice(sortedRules, func(i, j int) bool {
		return sortedRules[i].Priority < sortedRules[j].Priority
	})

	// 添加自定义规则
	for _, rule := range sortedRules {
		if !rule.Enabled {
			continue
		}

		routeRule := RouteRule{
			"outbound": rule.Outbound,
		}

		switch rule.RuleType {
		case "domain_suffix":
			routeRule["domain_suffix"] = rule.Values
		case "domain_keyword":
			routeRule["domain_keyword"] = rule.Values
		case "domain":
			routeRule["domain"] = rule.Values
		case "ip_cidr":
			routeRule["ip_cidr"] = rule.Values
		case "port":
			routeRule["port"] = rule.Values
		case "geosite":
			var tags []string
			for _, v := range rule.Values {
				tags = append(tags, fmt.Sprintf("geosite-%s", v))
			}
			routeRule["rule_set"] = tags
		case "geoip":
			var tags []string
			for _, v := range rule.Values {
				tags = append(tags, fmt.Sprintf("geoip-%s", v))
			}
			routeRule["rule_set"] = tags
		}

		rules = append(rules, routeRule)
	}

	// 添加规则组的路由规则
	for _, rg := range b.ruleGroups {
		if !rg.Enabled {
			continue
		}

		// Site 规则
		if len(rg.SiteRules) > 0 {
			var tags []string
			for _, sr := range rg.SiteRules {
				tags = append(tags, fmt.Sprintf("geosite-%s", sr))
			}
			rules = append(rules, RouteRule{
				"rule_set": tags,
				"outbound": rg.Name,
			})
		}

		// IP 规则
		if len(rg.IPRules) > 0 {
			var tags []string
			for _, ir := range rg.IPRules {
				tags = append(tags, fmt.Sprintf("geoip-%s", ir))
			}
			rules = append(rules, RouteRule{
				"rule_set": tags,
				"outbound": rg.Name,
			})
		}
	}

	route.Rules = rules

	return route
}

// buildExperimental 构建实验性配置
func (b *ConfigBuilder) buildExperimental() *ExperimentalConfig {
	return &ExperimentalConfig{
		ClashAPI: &ClashAPIConfig{
			ExternalController:    fmt.Sprintf("127.0.0.1:%d", b.settings.ClashAPIPort),
			ExternalUI:            b.settings.ClashUIPath,
			ExternalUIDownloadURL: "https://github.com/Zephyruso/zashboard/archive/refs/heads/gh-pages.zip",
			DefaultMode:           "rule",
		},
		CacheFile: &CacheFileConfig{
			Enabled: true,
			Path:    "cache.db",
		},
	}
}
