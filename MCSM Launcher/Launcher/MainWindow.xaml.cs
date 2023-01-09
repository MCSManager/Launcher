using ModernWpf.Controls;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Text.Json.Nodes;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;


namespace Launcher
{
    /// <summary>
    /// MainWindow.xaml 的交互逻辑
    /// </summary>
    public partial class MainWindow : Window
    {
        public MainWindow()
        {
            InitializeComponent();
        }

        private static Hardcodet.Wpf.TaskbarNotification.TaskbarIcon TaskbarIcon { get; set; } = new()
        {
            IconSource = new BitmapImage(new Uri("pack://application:,,,/Launcher;component/icon.ico")),
            ToolTipText = "MCSM已挂至后台，单击显示",
        };

        protected override async void OnInitialized(EventArgs e)
        {
            string file = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "web", "data", "SystemConfig", "config.json");
            if (!File.Exists(file))
            {
                MessageBox.Show("运行库不完整,请重新下载","MCSM Launcher",MessageBoxButton.OK,MessageBoxImage.Error);
                App.Current.Shutdown();
                
            }
            using var reader = new StreamReader(file);
            var node = JsonNode.Parse(await reader.ReadToEndAsync());
            MS_Run_State.Text = "关闭";
            MS_Run_Port.Text = node["httpPort"].ToString();
            UrltoPanel.NavigateUri = new($"http://localhost:{node["httpPort"]}");


            var menui1 = new MenuItem() { Header = "退出" };
            menui1.Click += (sender, args) =>
            {
                processes.ForEach(process =>
                {
                    if (!process.HasExited) process.Kill();
                });
                App.Current.Shutdown();
            };

            TaskbarIcon.ContextMenu = new ContextMenu()
            {
                Items =
                {
                    menui1
                }
            };
            TaskbarIcon.TrayLeftMouseUp += (sender, args) =>
            {
                Activate();
                Visibility = Visibility.Visible;
            };

            base.OnInitialized(e);
        }
        private static List<Process> processes { get; set; } = new();

        protected override void OnClosing(CancelEventArgs e)
        {
            e.Cancel = true;
            Visibility = Visibility.Hidden;
        }

        private void ToggleSwitch_Toggled(object sender, RoutedEventArgs e)
        {
            if (sender is ToggleSwitch toggle)
            {
                
                string file = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "node_app.exe");
                if (toggle.IsOn)
                {
                    var p1 = Process.Start(new ProcessStartInfo(file, $"\"{Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "daemon", "app.js")}\"")
                    {
                        CreateNoWindow = true,
                        UseShellExecute = false,
                        WorkingDirectory = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager","daemon")
                    });
                    var p2 = Process.Start(new ProcessStartInfo(file, $"\"{Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "web", "app.js")}\"")
                    {
                        CreateNoWindow = true,
                        UseShellExecute = false,
                        WorkingDirectory = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "web")
                    });
                    p1.EnableRaisingEvents = p2.EnableRaisingEvents = true;
                    if (!p1.HasExited && !p2.HasExited)
                    {
                        processes.Add(p1);
                        processes.Add(p2);
                        p1.Exited += OnExited;
                        p2.Exited += OnExited;
                        void OnExited(object sender,EventArgs e)
                        {
                            App.Current.Dispatcher.Invoke(() => MS_Run_State.Text = "关闭");
                            if (!p1.HasExited) p1.Kill();
                            if (!p2.HasExited) p2.Kill();
                            processes.Clear();
                        };
                        MS_Run_State.Text = "开启";
                        return;
                    }
                    if (!p1.HasExited) p1.Kill();
                    if (!p2.HasExited) p2.Kill();
                    

                }
                else
                {

                    processes.ForEach(process =>
                    {
                        if (!process.HasExited) process.Kill();
                    });
                    MS_Run_State.Text = "关闭";
                }
            }
        }



        private void HyperlinkButton_Click_1(object sender, RoutedEventArgs e)
        {
            Visibility = Visibility.Hidden;
        }

        private void HyperlinkButton_PreviewMouseLeftButtonUp(object sender, MouseButtonEventArgs e)
        {
            try
            {
                Process.Start(Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "daemon", "logs", "current.log"));
            }
            catch { }
            try
            {
                Process.Start(Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "web", "logs", "current.log"));
            }
            catch { }
        }
    }
}
