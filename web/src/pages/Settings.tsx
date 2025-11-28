import { useEffect, useState } from 'react';
import { Card, CardBody, CardHeader, Input, Button, Switch, Chip } from '@nextui-org/react';
import { Save, Download, Upload, Terminal } from 'lucide-react';
import { useStore } from '../store';
import type { Settings as SettingsType } from '../store';
import { launchdApi } from '../api';

export default function Settings() {
  const { settings, fetchSettings, updateSettings } = useStore();
  const [formData, setFormData] = useState<SettingsType | null>(null);
  const [launchdStatus, setLaunchdStatus] = useState<{ installed: boolean; running: boolean } | null>(null);

  useEffect(() => {
    fetchSettings();
    fetchLaunchdStatus();
  }, []);

  useEffect(() => {
    if (settings) {
      setFormData(settings);
    }
  }, [settings]);

  const fetchLaunchdStatus = async () => {
    try {
      const res = await launchdApi.status();
      setLaunchdStatus(res.data.data);
    } catch (error) {
      console.error('获取 launchd 状态失败:', error);
    }
  };

  const handleSave = async () => {
    if (formData) {
      await updateSettings(formData);
    }
  };

  const handleInstallLaunchd = async () => {
    try {
      await launchdApi.install();
      await fetchLaunchdStatus();
    } catch (error) {
      console.error('安装 launchd 服务失败:', error);
    }
  };

  const handleUninstallLaunchd = async () => {
    if (confirm('确定要卸载 launchd 服务吗？')) {
      try {
        await launchdApi.uninstall();
        await fetchLaunchdStatus();
      } catch (error) {
        console.error('卸载 launchd 服务失败:', error);
      }
    }
  };

  if (!formData) {
    return <div>加载中...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">设置</h1>
        <Button
          color="primary"
          startContent={<Save className="w-4 h-4" />}
          onPress={handleSave}
        >
          保存设置
        </Button>
      </div>

      {/* sing-box 配置 */}
      <Card>
        <CardHeader>
          <Terminal className="w-5 h-5 mr-2" />
          <h2 className="text-lg font-semibold">sing-box 配置</h2>
        </CardHeader>
        <CardBody className="space-y-4">
          <Input
            label="sing-box 路径"
            placeholder="/usr/local/bin/sing-box"
            value={formData.singbox_path}
            onChange={(e) => setFormData({ ...formData, singbox_path: e.target.value })}
          />
          <Input
            label="配置文件路径"
            placeholder="data/generated/config.json"
            value={formData.config_path}
            onChange={(e) => setFormData({ ...formData, config_path: e.target.value })}
          />
        </CardBody>
      </Card>

      {/* 入站配置 */}
      <Card>
        <CardHeader>
          <Download className="w-5 h-5 mr-2" />
          <h2 className="text-lg font-semibold">入站配置</h2>
        </CardHeader>
        <CardBody className="space-y-4">
          <Input
            type="number"
            label="混合代理端口"
            placeholder="2080"
            value={String(formData.mixed_port)}
            onChange={(e) => setFormData({ ...formData, mixed_port: parseInt(e.target.value) || 2080 })}
          />
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">TUN 模式</p>
              <p className="text-sm text-gray-500">启用 TUN 模式进行透明代理</p>
            </div>
            <Switch
              isSelected={formData.tun_enabled}
              onValueChange={(enabled) => setFormData({ ...formData, tun_enabled: enabled })}
            />
          </div>
        </CardBody>
      </Card>

      {/* DNS 配置 */}
      <Card>
        <CardHeader>
          <Upload className="w-5 h-5 mr-2" />
          <h2 className="text-lg font-semibold">DNS 配置</h2>
        </CardHeader>
        <CardBody className="space-y-4">
          <Input
            label="代理 DNS"
            placeholder="https://1.1.1.1/dns-query"
            value={formData.proxy_dns}
            onChange={(e) => setFormData({ ...formData, proxy_dns: e.target.value })}
          />
          <Input
            label="直连 DNS"
            placeholder="https://dns.alidns.com/dns-query"
            value={formData.direct_dns}
            onChange={(e) => setFormData({ ...formData, direct_dns: e.target.value })}
          />
        </CardBody>
      </Card>

      {/* 控制面板配置 */}
      <Card>
        <CardHeader>
          <h2 className="text-lg font-semibold">控制面板</h2>
        </CardHeader>
        <CardBody className="space-y-4">
          <Input
            type="number"
            label="Web 管理端口"
            placeholder="9090"
            value={String(formData.web_port)}
            onChange={(e) => setFormData({ ...formData, web_port: parseInt(e.target.value) || 9090 })}
          />
          <Input
            type="number"
            label="Clash API 端口"
            placeholder="9091"
            value={String(formData.clash_api_port)}
            onChange={(e) => setFormData({ ...formData, clash_api_port: parseInt(e.target.value) || 9091 })}
          />
          <Input
            label="漏网规则出站"
            placeholder="Proxy"
            value={formData.final_outbound}
            onChange={(e) => setFormData({ ...formData, final_outbound: e.target.value })}
          />
        </CardBody>
      </Card>

      {/* 自动化设置 */}
      <Card>
        <CardHeader>
          <h2 className="text-lg font-semibold">自动化</h2>
        </CardHeader>
        <CardBody className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">配置变更后自动应用</p>
              <p className="text-sm text-gray-500">订阅刷新或规则变更后自动重启 sing-box</p>
            </div>
            <Switch
              isSelected={formData.auto_apply}
              onValueChange={(enabled) => setFormData({ ...formData, auto_apply: enabled })}
            />
          </div>
          <Input
            type="number"
            label="订阅自动更新间隔 (分钟)"
            placeholder="60"
            description="设置为 0 表示禁用自动更新"
            value={String(formData.subscription_interval)}
            onChange={(e) => setFormData({ ...formData, subscription_interval: parseInt(e.target.value) || 0 })}
          />
        </CardBody>
      </Card>

      {/* launchd 服务管理 */}
      <Card>
        <CardHeader className="flex justify-between items-center">
          <h2 className="text-lg font-semibold">系统服务 (launchd)</h2>
          {launchdStatus && (
            <Chip
              color={launchdStatus.installed ? 'success' : 'default'}
              variant="flat"
              size="sm"
            >
              {launchdStatus.installed ? '已安装' : '未安装'}
            </Chip>
          )}
        </CardHeader>
        <CardBody>
          <p className="text-sm text-gray-500 mb-4">
            安装 launchd 服务后，sing-box 将在系统启动时自动运行，并在崩溃后自动重启。
          </p>
          <div className="flex gap-2">
            {launchdStatus?.installed ? (
              <Button
                color="danger"
                variant="flat"
                onPress={handleUninstallLaunchd}
              >
                卸载服务
              </Button>
            ) : (
              <Button
                color="primary"
                onPress={handleInstallLaunchd}
              >
                安装服务
              </Button>
            )}
          </div>
        </CardBody>
      </Card>
    </div>
  );
}
