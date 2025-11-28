import { useEffect } from 'react';
import { Card, CardBody, CardHeader, Button, Chip } from '@nextui-org/react';
import { Play, Square, RefreshCw, Cpu, HardDrive, Wifi } from 'lucide-react';
import { useStore } from '../store';
import { serviceApi, configApi } from '../api';

export default function Dashboard() {
  const { serviceStatus, subscriptions, fetchServiceStatus, fetchSubscriptions } = useStore();

  useEffect(() => {
    fetchServiceStatus();
    fetchSubscriptions();

    // 每 5 秒刷新状态
    const interval = setInterval(fetchServiceStatus, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleStart = async () => {
    try {
      await serviceApi.start();
      await fetchServiceStatus();
    } catch (error) {
      console.error('启动失败:', error);
    }
  };

  const handleStop = async () => {
    try {
      await serviceApi.stop();
      await fetchServiceStatus();
    } catch (error) {
      console.error('停止失败:', error);
    }
  };

  const handleRestart = async () => {
    try {
      await serviceApi.restart();
      await fetchServiceStatus();
    } catch (error) {
      console.error('重启失败:', error);
    }
  };

  const handleApplyConfig = async () => {
    try {
      await configApi.apply();
      await fetchServiceStatus();
    } catch (error) {
      console.error('应用配置失败:', error);
    }
  };

  const totalNodes = subscriptions.reduce((sum, sub) => sum + sub.node_count, 0);
  const enabledSubs = subscriptions.filter(sub => sub.enabled).length;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800 dark:text-white">仪表盘</h1>

      {/* 服务状态卡片 */}
      <Card>
        <CardHeader className="flex justify-between items-center">
          <div className="flex items-center gap-3">
            <h2 className="text-lg font-semibold">sing-box 服务</h2>
            <Chip
              color={serviceStatus?.running ? 'success' : 'danger'}
              variant="flat"
              size="sm"
            >
              {serviceStatus?.running ? '运行中' : '已停止'}
            </Chip>
          </div>
          <div className="flex gap-2">
            {serviceStatus?.running ? (
              <>
                <Button
                  size="sm"
                  color="danger"
                  variant="flat"
                  startContent={<Square className="w-4 h-4" />}
                  onPress={handleStop}
                >
                  停止
                </Button>
                <Button
                  size="sm"
                  color="primary"
                  variant="flat"
                  startContent={<RefreshCw className="w-4 h-4" />}
                  onPress={handleRestart}
                >
                  重启
                </Button>
              </>
            ) : (
              <Button
                size="sm"
                color="success"
                startContent={<Play className="w-4 h-4" />}
                onPress={handleStart}
              >
                启动
              </Button>
            )}
            <Button
              size="sm"
              color="primary"
              onPress={handleApplyConfig}
            >
              应用配置
            </Button>
          </div>
        </CardHeader>
        <CardBody>
          <div className="grid grid-cols-3 gap-4">
            <div>
              <p className="text-sm text-gray-500">版本</p>
              <p className="font-medium">{serviceStatus?.version || '-'}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">进程 ID</p>
              <p className="font-medium">{serviceStatus?.pid || '-'}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">状态</p>
              <p className="font-medium">
                {serviceStatus?.running ? '正常运行' : '未运行'}
              </p>
            </div>
          </div>
        </CardBody>
      </Card>

      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardBody className="flex flex-row items-center gap-4">
            <div className="p-3 bg-blue-100 dark:bg-blue-900 rounded-lg">
              <Wifi className="w-6 h-6 text-blue-600 dark:text-blue-300" />
            </div>
            <div>
              <p className="text-sm text-gray-500">订阅数量</p>
              <p className="text-2xl font-bold">{enabledSubs} / {subscriptions.length}</p>
            </div>
          </CardBody>
        </Card>

        <Card>
          <CardBody className="flex flex-row items-center gap-4">
            <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
              <HardDrive className="w-6 h-6 text-green-600 dark:text-green-300" />
            </div>
            <div>
              <p className="text-sm text-gray-500">节点总数</p>
              <p className="text-2xl font-bold">{totalNodes}</p>
            </div>
          </CardBody>
        </Card>

        <Card>
          <CardBody className="flex flex-row items-center gap-4">
            <div className="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg">
              <Cpu className="w-6 h-6 text-purple-600 dark:text-purple-300" />
            </div>
            <div>
              <p className="text-sm text-gray-500">系统状态</p>
              <p className="text-2xl font-bold">正常</p>
            </div>
          </CardBody>
        </Card>
      </div>

      {/* 订阅列表预览 */}
      <Card>
        <CardHeader>
          <h2 className="text-lg font-semibold">订阅概览</h2>
        </CardHeader>
        <CardBody>
          {subscriptions.length === 0 ? (
            <p className="text-gray-500 text-center py-4">暂无订阅，请前往节点页面添加</p>
          ) : (
            <div className="space-y-3">
              {subscriptions.map((sub) => (
                <div
                  key={sub.id}
                  className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg"
                >
                  <div className="flex items-center gap-3">
                    <Chip
                      size="sm"
                      color={sub.enabled ? 'success' : 'default'}
                      variant="dot"
                    >
                      {sub.name}
                    </Chip>
                    <span className="text-sm text-gray-500">
                      {sub.node_count} 个节点
                    </span>
                  </div>
                  <span className="text-sm text-gray-400">
                    更新于 {new Date(sub.updated_at).toLocaleString()}
                  </span>
                </div>
              ))}
            </div>
          )}
        </CardBody>
      </Card>
    </div>
  );
}
